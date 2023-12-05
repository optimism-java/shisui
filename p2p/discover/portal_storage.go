package discover

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"path"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/holiman/uint256"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ContentNotFound = fmt.Errorf("content not found")

	maxDistance = uint256.MustFromHex("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
)

const (
	sqliteName = "shisui.sqlite"
)

type Storage interface {
	ContentId(contentKey []byte) []byte

	Get(contentKey []byte, contentId []byte) ([]byte, error)

	Put(contentKey []byte, content []byte) error
}

func getDistance(a *uint256.Int, b *uint256.Int) *uint256.Int {
	return new(uint256.Int).Xor(a, b)
}

type PortalStorage struct {
	nodeId                 enode.ID
	nodeDataDir            string
	storageCapacityInBytes uint64
	radius                 *uint256.Int
	sqliteDB               *sql.DB
	getStmt                *sql.Stmt
	putStmt                *sql.Stmt
	delStmt                *sql.Stmt
	containStmt            *sql.Stmt
}

func NewPortalStorage(storageCapacityInBytes uint64, nodeId enode.ID, nodeDataDir string) (*PortalStorage, error) {

	sqlDb, err := sql.Open("sqlite3", path.Join(nodeDataDir, sqliteName))

	if err != nil {
		return nil, err
	}
	portalStorage := &PortalStorage{
		nodeId:                 nodeId,
		nodeDataDir:            nodeDataDir,
		storageCapacityInBytes: storageCapacityInBytes,
		radius:                 maxDistance,
		sqliteDB:               sqlDb,
	}

	err = portalStorage.createTable()
	if err != nil {
		return nil, err
	}

	err = portalStorage.initStmts()

	// Check whether we already have data, and use it to set radius

	return portalStorage, err
}

func (p *PortalStorage) ContentId(contentKey []byte) []byte {
	digest := sha256.Sum256(contentKey)
	return digest[:]
}

func (p *PortalStorage) Get(contentKey []byte, contentId []byte) ([]byte, error) {
	var res []byte
	err := p.getStmt.QueryRow().Scan(&res)
	if err == sql.ErrNoRows {
		return nil, ContentNotFound
	}
	return res, err
}

func (p *PortalStorage) Put(contentKey []byte, content []byte) error {
	contentId := p.ContentId(contentKey)
	_, err := p.putStmt.Exec(contentId, content)
	return err
}

func (p *PortalStorage) Close() error {
	p.getStmt.Close()
	p.putStmt.Close()
	p.delStmt.Close()
	p.containStmt.Close()
	return p.sqliteDB.Close()
}

func (p *PortalStorage) createTable() error {
	stat, err := p.sqliteDB.Prepare(createSql)
	if err != nil {
		return err
	}
	defer stat.Close()
	_, err = stat.Exec()
	return err
}

func (p *PortalStorage) initStmts() error {
	var stat *sql.Stmt
	var err error
	if stat, err = p.sqliteDB.Prepare(getSql); err != nil {
		return nil
	}
	p.getStmt = stat
	if stat, err = p.sqliteDB.Prepare(putSql); err != nil {
		return nil
	}
	p.putStmt = stat
	if stat, err = p.sqliteDB.Prepare(deleteSql); err != nil {
		return nil
	}
	p.delStmt = stat
	if stat, err = p.sqliteDB.Prepare(containSql); err != nil {
		return nil
	}
	p.containStmt = stat
	return nil
}

// get database size, content size and similar
func (p *PortalStorage) Size() int64 {
	return 0
}

// SQLite Statements

const createSql = `CREATE TABLE IF NOT EXISTS kvstore (
	key BLOB PRIMARY KEY,
	value BLOB
);`
const getSql = "SELECT value FROM kvstore WHERE key = (?1);"
const putSql = "INSERT OR REPLACE INTO kvstore (key, value) VALUES (?1, ?2);"
const deleteSql = "DELETE FROM kvstore WHERE key = (?1);"
const clearSql = "DELETE FROM kvstore"
const containSql = "SELECT 1 FROM kvstore WHERE key = (?1);"

const XOR_FIND_FARTHEST_QUERY = `SELECT
		content_id_long
		FROM content_metadata
		ORDER BY ((?1 | content_id_short) - (?1 & content_id_short)) DESC`

const CONTENT_KEY_LOOKUP_QUERY = "SELECT content_key FROM content_metadata WHERE content_id_long = (?1)"

const TOTAL_DATA_SIZE_QUERY = "SELECT TOTAL(content_size) FROM content_metadata"

const TOTAL_ENTRY_COUNT_QUERY = "SELECT COUNT(content_id_long) FROM content_metadata"

const PAGINATE_QUERY = "SELECT content_key FROM content_metadata ORDER BY content_key LIMIT :limit OFFSET :offset"

const CONTENT_SIZE_LOOKUP_QUERY = "SELECT content_size FROM content_metadata WHERE content_id_long = (?1)"
