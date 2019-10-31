package database

import (
	"github.com/jackc/pgx"
)

// Database ...
type Database struct {
	Conn *pgx.ConnPool
}

// CreateConnection ...
func CreateConnection(psqlURI string) (*Database, error) {
	pgxConfig, err := pgx.ParseURI(psqlURI)
	if err != nil {
		return nil, err
	}

	var db Database

	if db.Conn, err = pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig: pgxConfig,
		}); err != nil {
		return nil, err
	}

	return &db, nil
}

// Disconnect ...
func (db *Database) Disconnect() {
	db.Conn.Close()
}
