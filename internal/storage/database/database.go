package database

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
)

// SQL запросы для дальнейших функций.
var (
	QuerryCreate = `CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
	user_id VARCHAR(16),
    hash_id 	VARCHAR(8),
    income_url VARCHAR(100) NOT NULL UNIQUE,
	delete_flag BOOLEAN NOT NULL DEFAULT FALSE
	);`
	QuerrySave = `INSERT INTO urls (id, user_id, hash_id, income_url)
	VALUES  (DEFAULT, $1, $2, $3);`
	QuerryGet = `SELECT income_url, delete_flag
	FROM urls WHERE hash_id = $1 LIMIT 1;`
	QuerryGetOrig = `SELECT hash_id
	FROM urls WHERE income_url = $1 LIMIT 1;`
	QuerryGetUser = `SELECT hash_id, income_url
	FROM urls WHERE user_id = $1;`
	QuerryDeleteURL = `UPDATE urls SET delete_flag=TRUE
	WHERE user_id = ($1) AND hash_id = ANY($2);`
	QuerryStats = `SELECT COUNT(*), COUNT( DISTINCT user_id)
	FROM urls;`

	QuerryDropTable = `DROP TABLE urls;` // только для тестов!!!
)

// DBStor структура для интерфейсов базы данных.
type DBStor struct {
	DB      *sql.DB // соединение с базой
	DBInfo  string  // информация для подключения к базе
	InFiles bool    // флаг, указывающий на характер хранения данных (true - хранение в файле)
}

// DBinit инициализация хранилища и создание/подключение к таблице.
func DBinit(DBInfo string) (*DBStor, error) {
	var db DBStor
	var err error

	db.DBInfo = DBInfo

	db.DB, err = sql.Open("pgx", DBInfo)
	if err != nil && DBInfo != "" {
		return &DBStor{InFiles: true}, err
	} else if DBInfo == "" {
		return &DBStor{InFiles: true}, errors.New("turning off data base mode by command dbinfo = _")
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

// Save сохранение полученных значений в таблицу SQL.
func (db *DBStor) Save(key string, val string, UserID string) error {
	if db.InFiles {
		return nil
	}

	var err error
	_, err = db.DB.Exec(QuerrySave, UserID, key, val)
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

// Get получение значений из таблицы SQL по ключу.
func (db *DBStor) Get(key string) (string, error) {
	if db.InFiles {
		return "", nil
	}
	var err error
	var val string
	var isDelete = false

	row := db.DB.QueryRow(QuerryGet, key)
	err = row.Scan(&val, &isDelete)
	if err != nil {
		return "", err
	}
	if isDelete {
		return "", models.ErrorDeleted
	}
	return val, nil
}

// GetOrig получение значений из таблицы SQL по значению
// (для случаевв если переданный URL сожержится в базе).
func (db *DBStor) GetOrig(val string) (string, error) {
	var err error
	var key string

	row := db.DB.QueryRow(QuerryGetOrig, val)
	err = row.Scan(&key)
	if err != nil {
		return "", err
	}
	return key, nil
}

// GetUser получение всех сокращенных и оригинальных URL для конкретного пользователя.
func (db *DBStor) GetUser(UserID string, BaseAdr string) (models.UserBuff, error) {
	if db.InFiles {
		return nil, nil
	}
	var err error
	var OutBuff models.UserBuff
	rows, err := db.DB.Query(QuerryGetUser, UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var UserURL models.UserURL
		err = rows.Scan(&UserURL.ShortURL, &UserURL.OrigURL)
		if err != nil {
			return nil, err
		}
		UserURL.ShortURL = BaseAdr + "/" + UserURL.ShortURL

		OutBuff = append(OutBuff, UserURL)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(OutBuff) == 0 {
		return nil, models.ErrorNoUserURL
	}

	return OutBuff, nil
}

// SaveTx сохранение значений в таблицу при помощи транзакций
// (для случая с большим количеством URL на входе).
func (db *DBStor) SaveTx(InURLs models.InBuff, BaseAdr string) (models.OutBuff, error) {
	if db.InFiles {
		return nil, nil
	}

	OutBuff := make(models.OutBuff, len(InURLs))

	//подготовка транзакции
	tx, err := db.DB.Begin()
	if err != nil {
		db.InFiles = true
		return nil, err
	}
	// потоковое чтение JSON и сохранение в базу по транзакциям.
	for i, income := range InURLs {
		hashStr := income.Hash
		user := income.UserID

		_, err = tx.Exec(QuerrySave, user, hashStr, income.URL)
		if err != nil {
			tx.Rollback()
			db.InFiles = true
			return nil, err
		}

		var OutURL models.BatchOutURL
		OutURL.ShortURL = BaseAdr + "/" + hashStr
		OutURL.CorrID = income.CorrID

		OutBuff[i] = OutURL
	}
	err = tx.Commit()
	if err != nil {
		db.InFiles = true
		return nil, err
	}

	return OutBuff, nil
}

// Delete проставление флагов в базу о том, что URL удален из базы данных.
func (db *DBStor) Delete(keys []string, UserID string) error {
	if db.InFiles {
		return nil
	}

	var err error
	// в данном случае используется batch update для массива keys
	_, err = db.DB.Exec(QuerryDeleteURL, UserID, keys)
	if err != nil {
		db.InFiles = true
		return err
	}
	return nil
}

// DropTable сброс таблицы (только для тестов!!!).
func (db *DBStor) DropTable() error {
	var err error
	_, err = db.DB.Exec(QuerryDropTable)
	if err != nil {
		return err
	}
	return nil
}

// Ping проверка соединения с базой данных.
func (db *DBStor) Ping() error {
	var err error
	if err = db.DB.Ping(); err != nil {
		return err
	}
	return nil
}

// CodeIsUniqueViolation проверка возвращаемой ошибки
// на ошибку уникальности.
func (db *DBStor) CodeIsUniqueViolation(err error) bool {
	strErr := err.Error()
	return strings.Contains(strErr, pgerrcode.UniqueViolation)
}

// CloseDB закрытие соединения с базой данных
func (db *DBStor) CloseDB() error {
	return db.DB.Close()
}

// Stats функция для получения количества сохраненых
// URL в базе и количества пользователей
func (db *DBStor) Stats() (models.URLStats, error) {

	var stats models.URLStats
	if db.InFiles {
		return stats, nil
	}

	row := db.DB.QueryRow(QuerryStats)
	err := row.Scan(&stats.NumUrls, &stats.NumUsers)
	if err != nil {
		return stats, err
	}
	return stats, nil
}
