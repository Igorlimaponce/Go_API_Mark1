package database

// internal/database/postgres.go encapsula a logica de conexao e interacao com o banco de dados Postgres.
import (
	"database/sql"
	"fmt"
)

type PostgresDB struct {
	Conn *sql.DB
}

func NewPostgresDB(ConnString string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", ConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{Conn: db}, nil
}

func (p *PostgresDB) Close() error {
	return p.Conn.Close()
}
