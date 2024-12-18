package goseed

import (
	"database/sql"
	"fmt"
	"strings"
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

func (p *PostgresAdapter) InsertRow(tableName string, columns []string, valuesList [][]any) error {
	// Prepare the batch insert
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ", tableName, strings.Join(columns, ", "))

	// Prepare placeholders
	var placeholderGroups []string
	var flattenedValues []any

	for i, values := range valuesList {
		placeholders := make([]string, len(values))
		for j := range values {
			placeholders[j] = fmt.Sprintf("$%d", i*len(values)+j+1)
		}

		placeholderGroups = append(placeholderGroups, fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")))
		flattenedValues = append(flattenedValues, values...)
	}

	query += strings.Join(placeholderGroups, ", ")

	// Execute the batch insert
	_, err := p.DB.Exec(query, flattenedValues...)
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
