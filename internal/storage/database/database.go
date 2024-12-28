package database

import (
	"database/sql"
	"fmt"

	"github.com/Arti9991/shortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var QuerryCreate = `CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    hash_id 	VARCHAR(8),
    income_url VARCHAR(100)
	);`
var QuerrySave = `INSERT INTO urls (id, hash_id, income_url)
	VALUES  (DEFAULT, $1, $2);`
var QuerryGet = `SELECT income_url
	FROM urls WHERE hash_id = $1 LIMIT 1;`

type DBStor struct {
	DB      *sql.DB
	DBInfo  string
	inFiles bool
}

func DBinit(DBInfo string) (*DBStor, error) {
	var db DBStor
	var err error

	db.DBInfo = DBInfo

	db.DB, err = sql.Open("pgx", DBInfo)
	if err != nil || DBInfo == "" {
		return &DBStor{inFiles: true}, err
	}
	defer db.DB.Close()
	if err = db.DB.Ping(); err != nil {
		return &DBStor{inFiles: true}, err
	}

	res, err := db.DB.Exec(QuerryCreate)
	if err != nil {
		return &DBStor{inFiles: true}, err
	}
	fmt.Println(res)
	logger.Log.Info("✓ connected to ShortURL db! ✓ Table created!")
	db.inFiles = false
	return &db, nil
}

func (db *DBStor) DBsave(key string, val string) error {
	if db.inFiles {
		return nil
	}

	var err error

	db.DB, err = sql.Open("pgx", db.DBInfo)
	if err != nil {
		db.inFiles = true
		return err
	}
	defer db.DB.Close()

	res, err := db.DB.Exec(QuerrySave, key, val)
	if err != nil {
		db.inFiles = true
		return err
	}
	fmt.Println(res)
	return nil
}

func (db *DBStor) DBget(key string) (string, error) {
	if db.inFiles {
		return "", nil
	}

	var err error
	var val string

	db.DB, err = sql.Open("pgx", db.DBInfo)
	if err != nil {
		db.inFiles = true
		return "", err
	}
	defer db.DB.Close()

	row := db.DB.QueryRow(QuerryGet, key)
	err = row.Scan(&val)
	if err != nil {
		return "", err
	}
	fmt.Println(val)
	return val, nil
}
