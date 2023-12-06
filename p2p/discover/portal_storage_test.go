package discover

import (
	"database/sql"
	"testing"
)

func getDB() *sql.DB {
	const createSql = `CREATE TABLE IF NOT EXISTS kvstore (
		key BLOB PRIMARY KEY,
		value BLOB
	);`
	sqlDb, _ := sql.Open("sqlite3", "test.db")
	stat, _ :=sqlDb.Prepare(createSql)
	defer stat.Close()
	_, _ = stat.Exec()
	return sqlDb
}

func BenchmarkPrepare(b *testing.B) {
	const sql = "SELECT value FROM kvstore"
	db := getDB()
	b.ResetTimer()
	for i:=0; i < b.N; i++ {
		db.Exec(sql)
	}
}

func BenchmarkNoPrepare(b *testing.B) {
	const sql = "SELECT value FROM kvstore"
	db := getDB()
	stmp, _ := db.Prepare(sql)
	b.ResetTimer()
	for i:=0; i < b.N; i++ {
		stmp.Exec()
	}
}