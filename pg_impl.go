package goseed

import (
	"database/sql"
)

type PostgresAdapter struct {
	DB *sql.DB
}

func (p *PostgresAdapter) Ping() error {
	return p.DB.Ping()
}

func (p *PostgresAdapter) IsTableExists(tableName string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)`
	var exists bool
	err := p.DB.QueryRow(query, tableName).Scan(&exists)
	return exists, err
}

func (p *PostgresAdapter) InsertRow(tableName string, columns []string, values []any) error {
	query := buildInsertQuery(tableName, columns)
	_, err := p.DB.Exec(query, values...)
	return err
}

func (p *PostgresAdapter) GetColumns(tableName string) ([]string, error) {
	query := `SELECT column_name FROM information_schema.columns WHERE table_name = $1`
	rows, err := p.DB.Query(query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}
	return columns, nil
}
