package database

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func DBinit(DBInfo string) (*sql.DB, error) {

	db, err := sql.Open("pgx", DBInfo)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("✓ connected to books db")

	return db, nil

}
