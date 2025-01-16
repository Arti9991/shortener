package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var QuerryCreate = `CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    hash_id 	VARCHAR(8),
    income_url VARCHAR(100) NOT NULL UNIQUE
	);`
var QuerrySave = `INSERT INTO urls (id, hash_id, income_url)
	VALUES  (DEFAULT, $1, $2);`
var QuerryGet = `SELECT income_url
	FROM urls WHERE hash_id = $1 LIMIT 1;`
var QuerryGetOrig = `SELECT hash_id
	FROM urls WHERE income_url = $1 LIMIT 1;`

type DBStor struct {
	DB      *sql.DB
	DBInfo  string
	InFiles bool // флаг, указывающий на характер хранения данных (true - хранение в файле)
}

// инициализация хранилища и создание/подключение к таблице
func DBinit(DBInfo string) (*DBStor, error) {
	var db DBStor
	var err error

	db.DBInfo = DBInfo

	db.DB, err = sql.Open("pgx", DBInfo)
	if err != nil || DBInfo == "" {
		return &DBStor{InFiles: true}, err
	}
	defer db.DB.Close()
	if err = db.DB.Ping(); err != nil {
		return &DBStor{InFiles: true}, err
	}

	_, err = db.DB.Exec(QuerryCreate)
	if err != nil {
		return &DBStor{InFiles: true}, err
	}
	logger.Log.Info("✓ connected to ShortURL db!")
	db.InFiles = false
	return &db, nil
}

// сохранение полученных значений в таблицу SQL
func (db *DBStor) DBsave(key string, val string) error {
	if db.InFiles {
		return nil
	}

	var err error

	db.DB, err = sql.Open("pgx", db.DBInfo)
	if err != nil {
		db.InFiles = true
		return err
	}
	defer db.DB.Close()

	_, err = db.DB.Exec(QuerrySave, key, val)
	if err != nil {
		if db.CodeIsUniqueViolation(err) {
			return err
		} else {
			db.InFiles = true
			return err
		}
	}
	return nil
}

// получение значений из таблицы SQL по ключу
func (db *DBStor) DBget(key string) (string, error) {
	if db.InFiles {
		return "", nil
	}

	var err error
	var val string

	db.DB, err = sql.Open("pgx", db.DBInfo)
	if err != nil {
		db.InFiles = true
		return "", err
	}
	defer db.DB.Close()

	row := db.DB.QueryRow(QuerryGet, key)
	err = row.Scan(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}

// получение значений из таблицы SQL по значению
// (для случаевв если переданный URL сожержится в базе)
func (db *DBStor) DBgetOrig(val string) (string, error) {
	var err error
	var key string

	db.DB, err = sql.Open("pgx", db.DBInfo)
	if err != nil {
		db.InFiles = true
		return "", err
	}
	defer db.DB.Close()

	row := db.DB.QueryRow(QuerryGetOrig, val)
	err = row.Scan(&key)
	if err != nil {
		return "", err
	}
	return key, nil
}

// сохранение значений в таблицу при помощи транзакций
// (для случая с большим количеством URL на входе)
func (db *DBStor) DBsaveTx(key string, val string) error {
	if db.InFiles {
		return nil
	}

	var err error

	db.DB, err = sql.Open("pgx", db.DBInfo)
	if err != nil {
		db.InFiles = true
		return err
	}
	defer db.DB.Close()

	tx, err := db.DB.Begin()
	if err != nil {
		db.InFiles = true
		return err
	}

	_, err = tx.Exec(QuerrySave, key, val)
	if err != nil {
		tx.Rollback()
		db.InFiles = true
		return err
	}
	return tx.Commit()
}

// проверка соединения с базой данных
func (db *DBStor) Ping() error {
	var err error
	db.DB, err = sql.Open("pgx", db.DBInfo)
	if err != nil {
		return err
	}
	defer db.DB.Close()
	if err = db.DB.Ping(); err != nil {
		return err
	}
	return nil
}

func (db *DBStor) CodeIsUniqueViolation(err error) bool {
	strErr := fmt.Sprintf("%s", err)
	arrErr := strings.Split(strErr, "(SQLSTATE")
	if len(arrErr) < 2 {
		return false
	}
	arrErr[1], _ = strings.CutSuffix(arrErr[1], ")")
	arrErr[1], _ = strings.CutPrefix(arrErr[1], " ")
	return arrErr[1] == pgerrcode.UniqueViolation
}
