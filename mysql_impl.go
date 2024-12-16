package goseed

import "database/sql"

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

func (m *MySQLAdapter) InsertRow(tableName string, columns []string, values []interface{}) error {
	query := buildInsertQuery(tableName, columns)
	_, err := m.DB.Exec(query, values...)
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
