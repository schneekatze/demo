package database

import "database/sql"

type SourceConnection interface {
	// Open connects to the database, returning either a db or error on opening the db handle
	// Important — it's just the pool of connections. Doesn't mean the connection itself is established
	// It's unclear yet, but the idea is you should not really use the handle directly
	// Rather than this interface.
	Open() error

	// Close closes the opened handle pool
	Close() error
}

type SQLDriverConnection interface {
	// Prepare is a proxy of sql.PrepareContext (adds our own context)
	Prepare(query string) (*sql.Stmt, error)

	// Exec is a proxy of sql.ExecContext (adds our own context)
	Exec(query string, args ...interface{}) (sql.Result, error)

	// Query is a proxy of sql.ExecQuery (adds our own context)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type ViaSshConnection interface {
	// OpenSsh Open connects to the database via SSH tunnel, returning either a db or error on opening the db handle
	// Important — it's just the pool of connections. Doesn't mean the connection itself is established
	OpenSsh(host string, port uint16, user string) error
}
