package database

import (
	"context"
	"database/sql"
	"fmt"
)

var mysql MySQLSourceConnection

type MySQLSourceConnection struct {
	Host     string
	Port     uint16
	User     string
	Password string
	Database string
	Context  context.Context
	DB       *sql.DB
}

func (s *MySQLSourceConnection) Open(ctx context.Context) error {
	var err error
	s.DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s:%d/%s", s.User, s.Password, s.Host, s.Port, s.Database))
	if err != nil {
		return err
	}

	s.Context = ctx

	return nil
}

func (s *MySQLSourceConnection) Close() error {
	err := s.DB.Close()

	return err
}

func (s *MySQLSourceConnection) Prepare(query string) (*sql.Stmt, error) {
	return s.DB.PrepareContext(s.Context, query)
}

func (s *MySQLSourceConnection) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.DB.ExecContext(s.Context, query, args...)
}

func (s *MySQLSourceConnection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.DB.QueryContext(s.Context, query, args...)
}
