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
	InFiles bool
}

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

	res, err := db.DB.Exec(QuerryCreate)
	if err != nil {
		return &DBStor{InFiles: true}, err
	}
	fmt.Println(res)
	logger.Log.Info("✓ connected to ShortURL db!")
	db.InFiles = false
	return &db, nil
}

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

	res, err := db.DB.Exec(QuerrySave, key, val)
	if err != nil {
		db.InFiles = true
		return err
	}
	fmt.Println(res)
	return nil
}

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
	fmt.Println(val)
	return val, nil
}

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

	res, err := tx.Exec(QuerrySave, key, val)
	if err != nil {
		tx.Rollback()
		db.InFiles = true
		return err
	}
	fmt.Println(res)
	return tx.Commit()
}

// func (db *DBStor) DBsaveMany(decoder *json.Decoder, hashStr string) ([]byte, error) {
// 	if db.InFiles {
// 		return nil, nil
// 	}
// 	var OutBuff []byte
// 	IncomeURL := &struct {
// 		Corr_id string `json:"correlation_id"`
// 		URL     string `json:"url"`
// 	}{}
// 	OutURL := &struct {
// 		Corr_id  string `json:"correlation_id"`
// 		ShortURL string `json:"short_url"`
// 	}{}

// 	stmt, err := db.DB.Prepare(QuerrySave)
// 	if err != nil {
// 		db.InFiles = true
// 		return nil, err
// 	}
// 	defer stmt.Close()

// 	for {
// 		err := decoder.Decode(&IncomeURL)
// 		if err == io.EOF {
// 			break
// 		} else if err != nil {
// 			return nil, err
// 		}

// 		// hd.dt.AddValue(hashStr, IncomeURL.URL)

// 		// err = hd.Files.FileSave(hashStr, IncomeURL.URL)
// 		// if err != nil {
// 		// 	logger.Log.Info("Error in FileSave", zap.Error(err))
// 		// }

// 		if !hd.DataBase.InFiles {
// 			_, err := stmt.Exec(hashStr, IncomeURL.URL)
// 			if err != nil {
// 				logger.Log.Info("Error in DB Save", zap.Error(err))
// 				hd.DataBase.InFiles = true
// 			}
// 		}

// 		OutURL.ShortURL = hd.BaseAdr + "/" + hashStr
// 		OutURL.Corr_id = IncomeURL.Corr_id
// 		tmp, err := json.Marshal(OutURL)
// 		if err != nil {
// 			logger.Log.Info("Wrong responce body", zap.Error(err))
// 			res.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		OutBuff = append(OutBuff, tmp...)
// 	}
// }
