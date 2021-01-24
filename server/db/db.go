package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Manager struct {
	DB *sql.DB
}
type BlockInfo struct {
	BlockNumber  int
	Transactions int
	Total        float64
}

func Connect(psqlInfo string) *sql.DB {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err.Error())
	}

	return db
}

func CreateTables(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS blocks (
		block_number integer UNIQUE NOT NULL,
		transactions integer NOT NULL,
		total double precision NOT NULL
	);`
	_, err := db.Exec(query)
	if err != nil {
		panic(err.Error())
	}
}

func (m *Manager) GetBlockInfo(blockNumber int) (error, bool, BlockInfo) {
	var blockInfo BlockInfo
	row := m.DB.QueryRow("SELECT transactions, total FROM blocks WHERE block_number = $1", blockNumber)
	switch err := row.Scan(&blockInfo.Transactions, &blockInfo.Total); err {
	case nil:
		return nil, true, blockInfo
	case sql.ErrNoRows:
		return nil, false, blockInfo
	default:
		return err, false, blockInfo
	}
}

func (m *Manager) InsertBlockInfo(blockNumber, transactionsNumber int, totalEth float64) error {
	_, err := m.DB.Exec("INSERT INTO blocks (block_number, transactions, total) values($1, $2, $3)", blockNumber, transactionsNumber, totalEth)
	return err
}
