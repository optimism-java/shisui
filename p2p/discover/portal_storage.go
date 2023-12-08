package discover

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"path"
	"strings"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/holiman/uint256"
	"github.com/mattn/go-sqlite3"
)

var (
	ContentNotFound = fmt.Errorf("content not found")

	maxDistance = uint256.MustFromHex("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
)

const (
	sqliteName              = "shisui.sqlite"
	contentDeletionFraction = 0.05 // 5% of the content will be deleted when the storage capacity is hit and radius gets adjusted.
)

type Storage interface {
	ContentId(contentKey []byte) []byte

	Get(contentKey []byte, contentId []byte) ([]byte, error)

	Put(contentKey []byte, content []byte) error
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
	log                    log.Logger
}

func xor(contentId, nodeId []byte) []byte {
	// length of contentId maybe not 32bytes
	padding := make([]byte, 32) 
	if len(contentId) != len(nodeId) {
		copy(padding, contentId)
	} else {
		padding = contentId
	}
	res := make([]byte, len(padding))
	for i, _ := range padding {
		res[i] = padding[i] ^ nodeId[i]
	}
	return res
}

func NewPortalStorage(storageCapacityInBytes uint64, nodeId enode.ID, nodeDataDir string) (*PortalStorage, error) {

	sql.Register("sqlite3_custom", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			if err := conn.RegisterFunc("xor", xor, false); err != nil {
				return err
			}
			return nil
		},
	})

	sqlDb, err := sql.Open("sqlite3_custom", path.Join(nodeDataDir, sqliteName))

	if err != nil {
		return nil, err
	}
	portalStorage := &PortalStorage{
		nodeId:                 nodeId,
		nodeDataDir:            nodeDataDir,
		storageCapacityInBytes: storageCapacityInBytes,
		radius:                 maxDistance,
		sqliteDB:               sqlDb,
		log:                    log.New("protocol_storage"),
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
	err := p.getStmt.QueryRow(contentId).Scan(&res)
	if err == sql.ErrNoRows {
		return nil, ContentNotFound
	}
	return res, err
}
type PutResult struct {
	err    error
	pruned bool
	count  int
}

func (p *PutResult) Err() error {
	return p.err
}

func (p *PutResult) Pruned() bool {
	return p.pruned
}

func (p *PutResult) PrunedCount() int {
	return p.count
}

func newPutResultWithErr(err error) PutResult {
	return PutResult{
		err: err,
	}
}

func (p *PortalStorage) Put(contentKey []byte, content []byte) PutResult {
	return p.put(p.ContentId(contentKey), content)
}

func (p *PortalStorage) put(contentId []byte, content []byte) PutResult {
	_, err := p.putStmt.Exec(contentId, content)
	if err != nil {
		return newPutResultWithErr(err)
	}

	dbSize, err := p.UsedSize()
	if err != nil {
		return newPutResultWithErr(err)
	}
	if dbSize > p.storageCapacityInBytes {
		count, err := p.deleteContentFraction(contentDeletionFraction)
		//
		if err != nil {
			log.Warn("failed to delete oversize item")
			return newPutResultWithErr(err)
		}
		return PutResult{pruned: true, count: count}
	}

	return PutResult{}
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
func (p *PortalStorage) Size() (uint64, error) {
	sql := "SELECT page_count * page_size as size FROM pragma_page_count(), pragma_page_size();"
	stmt, err := p.sqliteDB.Prepare(sql)
	if err != nil {
		return 0, err
	}
	var res uint64
	err = stmt.QueryRow().Scan(&res)
	return res, err
}

func (p *PortalStorage) UnusedSize() (uint64, error) {
	sql := "SELECT freelist_count * page_size as size FROM pragma_freelist_count(), pragma_page_size();"
	return p.queryRowUint64(sql)
}

func (p *PortalStorage) UsedSize() (uint64, error) {
	size, err := p.Size()
	if err != nil {
		return 0, err
	}
	unusedSize, err := p.UnusedSize()
	if err != nil {
		return 0, err
	}
	return size - unusedSize, err
}

func (p *PortalStorage) ContentCount() (uint64, error) {
	sql := "SELECT COUNT(key) FROM kvstore;"
	return p.queryRowUint64(sql)
}

func (p *PortalStorage) ContentSize() (uint64, error) {
	sql := "SELECT SUM(length(value)) FROM kvstore"
	return p.queryRowUint64(sql)
}

func (p *PortalStorage) queryRowUint64(sql string) (uint64, error) {
	// sql := "SELECT SUM(length(value)) FROM kvstore"
	stmt, err := p.sqliteDB.Prepare(sql)
	if err != nil {
		return 0, err
	}
	var res uint64
	err = stmt.QueryRow().Scan(&res)
	return res, err
}

func (p *PortalStorage) GetLargestDistance() (*uint256.Int, error) {
	stmt, err := p.sqliteDB.Prepare(XOR_FIND_FARTHEST_QUERY)
	if err != nil {
		return nil, err
	}
	var distance []byte

	err = stmt.QueryRow(p.nodeId[:]).Scan(&distance)
	if err != nil {
		return nil, err
	}
	// reverse the distance, because big.SetBytes is big-endian
	reverseBytes(distance)
	// for i:= 0; i < len(distance) / 2; i++ {
	// 	distance[i], distance[len(distance) - i - 1] = distance[len(distance) - i - 1], distance[i]
	// }
	bigNum := new(big.Int).SetBytes(distance)
	res := uint256.MustFromBig(bigNum)
	return res, nil
}

func (p *PortalStorage) EstimateNewRadius(currentRadius *uint256.Int) (*uint256.Int, error) {
	currrentSize, err := p.UsedSize()
	if err != nil {
		return nil, err
	}
	sizeRatio := currrentSize / p.storageCapacityInBytes
	if sizeRatio > 0 {
		bigFormat := new(big.Int).SetUint64(sizeRatio)
		return new(uint256.Int).Div(currentRadius, uint256.MustFromBig(bigFormat)), nil
	}
	return currentRadius, nil
}

func (p *PortalStorage) deleteContentFraction(fraction float64) (deleteCount int, err error) {
	if fraction <= 0 || fraction >= 1 {
		return deleteCount, errors.New("fraction should be between 0 and 1")
	}
	totalContentSize, err := p.ContentSize()
	if err != nil {
		return deleteCount, err
	}
	bytesToDelete := uint64(fraction * float64(totalContentSize))
	// deleteElements := 0
	deleteBytes := 0

	rows, err := p.sqliteDB.Query(getAllOrderedByDistanceSql, p.nodeId[:])
	if err != nil {
		return deleteCount, err
	}
	defer rows.Close()
	idsToDelete := make([][]byte, 0)
	for deleteBytes < int(bytesToDelete) && rows.Next() {
		var contentId []byte
		var payloadLen int
		var distance []byte
		err = rows.Scan(&contentId, &payloadLen, &distance)
		if err != nil {
			return deleteCount, err
		}
		idsToDelete = append(idsToDelete, contentId)
		// err = p.del(contentId)
		if err != nil {
			return deleteCount, err
		}
		deleteBytes += payloadLen
		deleteCount++
	}
	// row must close first, or database is locked
	// rows.Close() can call multi times
	rows.Close()
	err = p.batchDel(idsToDelete)
	return
}

func (p *PortalStorage) del(contentId []byte) error {
	_, err := p.delStmt.Exec(contentId)
	return err
}

func (p *PortalStorage) batchDel(ids [][]byte) error {
	query := "DELETE FROM kvstore WHERE key IN (?" + strings.Repeat(", ?", len(ids)-1) + ")"
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	// 执行删除操作
	_, err := p.sqliteDB.Exec(query, args...)
	return err
}

func (p *PortalStorage) ReclaimSpace() error {
	_, err := p.sqliteDB.Exec("VACUUM;")
	return err
}

func (p *PortalStorage) deleteContentOutOfRadius(radius *uint256.Int) error {
	res, err := p.sqliteDB.Exec(deleteOutOfRadiusStmt, p.nodeId[:], radius)
	count, _ :=  res.RowsAffected()
	p.log.Trace("delete %d items",count)
	return err
}

func (p *PortalStorage) ForcePrune(radius *uint256.Int) error {
	return p.deleteContentOutOfRadius(radius)
}

func reverseBytes(src []byte) {
	for i:= 0; i < len(src) / 2; i++ {
		src[i], src[len(src) - i - 1] = src[len(src) - i - 1], src[i]
	}
}

// SQLite Statements

const createSql = `CREATE TABLE IF NOT EXISTS kvstore (
	key BLOB PRIMARY KEY,
	value BLOB
);`
const getSql = "SELECT value FROM kvstore WHERE key = (?1);"
const putSql = "INSERT OR REPLACE INTO kvstore (key, value) VALUES (?1, ?2);"
const deleteSql = "DELETE FROM kvstore WHERE key = (?1);"

// const clearSql = "DELETE FROM kvstore"
const containSql = "SELECT 1 FROM kvstore WHERE key = (?1);"
const getAllOrderedByDistanceSql = "SELECT key, length(value), xor(key, (?1)) as distance FROM kvstore ORDER BY distance DESC;"
const deleteOutOfRadiusStmt = "DELETE FROM kvstore WHERE xor(key, (?1)) < ?2"

const XOR_FIND_FARTHEST_QUERY = `SELECT
		xor(key, (?1)) as distance
		FROM kvstore
		ORDER BY distance DESC`

