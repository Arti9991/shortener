package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/exp/rand"
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

	row := db.DB.QueryRow(QuerryGetOrig, val)
	err = row.Scan(&key)
	if err != nil {
		return "", err
	}
	return key, nil
}

// сохранение значений в таблицу при помощи транзакций
// (для случая с большим количеством URL на входе)
func (db *DBStor) DBsaveTx(dec *json.Decoder, BaseAdr string) (models.OutBuff, error) {
	if db.InFiles {
		return nil, nil
	}
	var IncomeURL models.BatchIncomeURL
	var OutBuff models.OutBuff

	tx, err := db.DB.Begin()
	if err != nil {
		db.InFiles = true
		return nil, err
	}
	for dec.More() {
		err := dec.Decode(&IncomeURL)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		hashStr := randomString(8)
		//inmemory.AddValue(hashStr, IncomeURL.URL)

		// err = hd.Files.FileSave(hashStr, IncomeURL.URL)
		// if err != nil {
		// 	logger.Log.Info("Error in FileSave", zap.Error(err))
		// }

		_, err = tx.Exec(QuerrySave, hashStr, IncomeURL.URL)
		if err != nil {
			tx.Rollback()
			db.InFiles = true
			return nil, err
		}

		var OutURL models.BatchOutURL
		OutURL.ShortURL = BaseAdr + "/" + hashStr
		OutURL.CorrID = IncomeURL.CorrID

		OutBuff = append(OutBuff, OutURL)
	}
	err = tx.Commit()
	if err != nil {
		db.InFiles = true
		return nil, err
	}

	// tx, err := db.DB.Begin()
	// if err != nil {
	// 	db.InFiles = true
	// 	return err
	// }

	// _, err = tx.Exec(QuerrySave, key, val)
	// if err != nil {
	// 	tx.Rollback()
	// 	db.InFiles = true
	// 	return err
	// }
	return OutBuff, nil
}

// проверка соединения с базой данных
func (db *DBStor) Ping() error {
	var err error
	defer db.DB.Close()
	if err = db.DB.Ping(); err != nil {
		return err
	}
	return nil
}

func (db *DBStor) CodeIsUniqueViolation(err error) bool {
	strErr := fmt.Sprintf("%s", err)
	return strings.Contains(strErr, pgerrcode.UniqueViolation)
}
func randomString(n int) string {

	var bt []byte
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for range n {
		bt = append(bt, charset[rand.Intn(len(charset))])
	}

	return string(bt)
}
