package discover

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"path"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/holiman/uint256"
	"github.com/linxGnu/grocksdb"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ContentNotFound = fmt.Errorf("content not found")

	maxDistance = uint256.MustFromHex("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
)

const (
	kvStoreName = "rocksdb"
	sqliteName  = "shisui.sqlite"
)

type Storage interface {
	ContentId(contentKey []byte) []byte

	Get(contentKey []byte, contentId []byte) ([]byte, error)

	Put(contentKey []byte, content []byte) error
}

func getDistance(a *uint256.Int, b *uint256.Int) *uint256.Int {
	return new(uint256.Int).Xor(a, b)
}

// opts := grocksdb.NewDefaultOptions()

// opts.SetCreateIfMissing(true)

// db, err := grocksdb.OpenDb(opts, "/path/to/db")

type PortalStorage struct {
	nodeId                 enode.ID
	nodeDataDir            string
	storageCapacityInBytes uint64
	radius                 *uint256.Int
	kv                     *grocksdb.DB
	sqliteDB               *sql.DB
}

func NewPortalStorage(storageCapacityInMB uint64, nodeId enode.ID, nodeDataDir string) (*PortalStorage, error) {

	opts := grocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)

	kv, err := grocksdb.OpenDb(opts, path.Join(nodeDataDir, kvStoreName))

	if err != nil {
		return nil, err
	}

	sqlDb, err := sql.Open("sqlite3", path.Join(nodeDataDir, sqliteName))

	if err != nil {
		return nil, err
	}

	portalStorage := &PortalStorage{
		nodeId:                 nodeId,
		nodeDataDir:            nodeDataDir,
		storageCapacityInBytes: storageCapacityInMB * 1000 * 1000,
		radius:                 maxDistance,
		kv:                     kv,
		sqliteDB:               sqlDb,
	}

	err = portalStorage.setupSql()

	if err != nil {
		return nil, err
	}

	// Check whether we already have data, and use it to set radius

	return portalStorage, nil
}

func (p *PortalStorage) ContentId(contentKey []byte) []byte {
	digest := sha256.Sum256(contentKey)
	return digest[:]
}

func (p *PortalStorage) Get(contentKey []byte, contentId []byte) ([]byte, error) {
	return p.kv.GetBytes(&grocksdb.ReadOptions{}, contentId)
}

func (p *PortalStorage) setupSql() error {
	stat, err := p.sqliteDB.Prepare(CREATE_QUERY)
	if err != nil {
		return err
	}
	defer stat.Close()
	_, err = stat.Exec()
	return err
}

func (p *PortalStorage) total_entry_count() (uint64, error) {
	stat, err := p.sqliteDB.Prepare(TOTAL_ENTRY_COUNT_QUERY)
	if err != nil {
		return 0, err
	}
	defer stat.Close()
	var total uint64
	err = stat.QueryRow().Scan(&total)
	return total, err
}

func (p *PortalStorage) capacityReached() (bool, error) {
	storageUsage, err := p.getTotalStorageUsageInBytesFromNetwork()
	return storageUsage > p.storageCapacityInBytes, err
}

// Internal method for measuring the total amount of requestable data that the node is storing.
func (p *PortalStorage) getTotalStorageUsageInBytesFromNetwork() (uint64, error) {
	stat, err := p.sqliteDB.Prepare(TOTAL_DATA_SIZE_QUERY)
	if err != nil {
		return 0, err
	}
	defer stat.Close()
	var total uint64
	err = stat.QueryRow().Scan(&total)
	return total, err
}

// func (p *PortalStorage) findFarthestContentId() ([]byte, error) {
// 	stat, err := p.sqliteDB.Prepare(XOR_FIND_FARTHEST_QUERY)
// 	if err != nil {
// 		return nil, err
// 	}
// }

// SQLite Statements
const CREATE_QUERY = `CREATE TABLE IF NOT EXISTS content_metadata (
	content_id_long TEXT PRIMARY KEY,
	content_id_short INTEGER NOT NULL,
	content_key TEXT NOT NULL,
	content_size INTEGER
);
CREATE INDEX content_size_idx ON content_metadata(content_size);
CREATE INDEX content_id_short_idx ON content_metadata(content_id_short);
CREATE INDEX content_id_long_idx ON content_metadata(content_id_long);`

const INSERT_QUERY = `INSERT OR IGNORE INTO content_metadata (content_id_long, content_id_short, content_key, content_size)
VALUES (?1, ?2, ?3, ?4)`

const DELETE_QUERY = `DELETE FROM content_metadata
WHERE content_id_long = (?1)`

const XOR_FIND_FARTHEST_QUERY = `SELECT
		content_id_long
		FROM content_metadata
		ORDER BY ((?1 | content_id_short) - (?1 & content_id_short)) DESC`

const CONTENT_KEY_LOOKUP_QUERY = "SELECT content_key FROM content_metadata WHERE content_id_long = (?1)"

const TOTAL_DATA_SIZE_QUERY = "SELECT TOTAL(content_size) FROM content_metadata"

const TOTAL_ENTRY_COUNT_QUERY = "SELECT COUNT(content_id_long) FROM content_metadata"

const PAGINATE_QUERY = "SELECT content_key FROM content_metadata ORDER BY content_key LIMIT :limit OFFSET :offset"

const CONTENT_SIZE_LOOKUP_QUERY = "SELECT content_size FROM content_metadata WHERE content_id_long = (?1)"
