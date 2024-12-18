package goseed

import (
	"database/sql"
	"fmt"
	"strings"
)

type MySQLAdapter struct {
	DB *sql.DB
}

func (m *MySQLAdapter) Ping() error {
	return m.DB.Ping()
}

func (m *MySQLAdapter) IsTableExists(tableName string) (bool, error) {
	query := `SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?`
	var count int
	err := m.DB.QueryRow(query, tableName).Scan(&count)
	return count > 0, err
}

func (m *MySQLAdapter) InsertRow(tableName string, columns []string, valuesList [][]any) error {
	// Prepare the batch insert query
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ", tableName, strings.Join(columns, ", "))

	// Prepare placeholders
	var placeholderGroups []string
	var flattenedValues []interface{}

	for _, values := range valuesList {
		placeholders := make([]string, len(values))
		for range values {
			placeholders = append(placeholders, "?")
		}
		placeholderGroups = append(placeholderGroups, fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")))
		flattenedValues = append(flattenedValues, values...)
	}

	query += strings.Join(placeholderGroups, ", ")

	// Execute the batch insert
	_, err := m.DB.Exec(query, flattenedValues...)
	return err
}

func (m *MySQLAdapter) GetColumns(tableName string) ([]string, error) {
	query := `SELECT column_name FROM information_schema.columns WHERE table_name = ?`
	rows, err := m.DB.Query(query, tableName)
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
