package database

import (
	"database/sql"
	"fmt"

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
	Db      *sql.DB
	DBInfo  string
	inFiles bool
}

func DBinit(DBInfo string) (*DBStor, error) {
	var db DBStor
	var err error

	db.DBInfo = DBInfo

	db.Db, err = sql.Open("pgx", DBInfo)
	if err != nil || DBInfo == "" {
		return &DBStor{inFiles: true}, err
	}
	defer db.Db.Close()
	if err = db.Db.Ping(); err != nil {
		return &DBStor{inFiles: true}, err
	}
	fmt.Println("✓ connected to ShortURL db")

	res, err := db.Db.Exec(QuerryCreate)
	if err != nil {
		return &DBStor{inFiles: true}, err
	}
	fmt.Println(res)
	fmt.Println("✓ Table created")
	db.inFiles = false
	return &db, nil

}

func (db *DBStor) DBsave(key string, val string) error {
	if db.inFiles {
		return nil
	}

	var err error

	db.Db, err = sql.Open("pgx", db.DBInfo)
	if err != nil {
		db.inFiles = true
		return err
	}
	defer db.Db.Close()

	res, err := db.Db.Exec(QuerrySave, key, val)
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

	db.Db, err = sql.Open("pgx", db.DBInfo)
	if err != nil {
		db.inFiles = true
		return "", err
	}
	defer db.Db.Close()

	row := db.Db.QueryRow(QuerryGet, key)
	err = row.Scan(&val)
	if err != nil {
		return "", err
	}
	fmt.Println(val)
	return val, nil
}
