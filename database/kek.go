package database

import (
	"challenge/config"
	"database/sql"
	"fmt"
)

func DbConn() *sql.DB {
	DBConfig := config.Cfg().DBConfig

	address := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		DBConfig.User,
		DBConfig.Password,
		DBConfig.Host,
		DBConfig.Port,
		DBConfig.Name,
	)
	db, err := sql.Open(
		"mysql",
		address,
	)

	if err != nil {
		panic(err)
	}

	return db
}
