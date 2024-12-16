package goseed

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"reflect"
	"strings"

	"github.com/go-faker/faker/v4"
)

/*
	This file contains an implementation of goseed, a database seeder library for PostgreSQL.
	Developed by: @ezrantn
	GitHub: github.com/ezrantn/goseed
	Contact: ezrantn@proton.me

	This library is licensed under the BSD-3 Clause. See LICENSE for more details.
*/

// Defines the structure for seeding a specific table
type TableSeeder struct {
	TableName string
	RowCount  int
	Model     interface{}
}

// Seeder orchestrates the database seeding process
type Seeder struct {
	DB           *sql.DB
	TableSeeders []TableSeeder
}

// Creates a new seeder instance
func NewGoSeed(db *sql.DB) (*Seeder, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("goseed unable to connect to database: %v", err)
	}

	return &Seeder{
		DB:           db,
		TableSeeders: []TableSeeder{},
	}, nil
}

// Adds a table seeder to the seeder instance
func (s *Seeder) Add(seeder TableSeeder) error {
	if seeder.TableName == "" {
		return errors.New("table name cannot be empty")
	}

	if seeder.Model == "" {
		return errors.New("model struct cannot be nil")
	}

	if seeder.RowCount <= 0 {
		return errors.New("row count must be greater than zero")
	}

	s.TableSeeders = append(s.TableSeeders, seeder)
	return nil
}

// Executes the seeding process for all configured table
func (s *Seeder) Run() error {
	for _, table := range s.TableSeeders {
		log.Printf("Seeding %d rows for table '%s'...\n", table.RowCount, table.TableName)

		// Check if the table exists
		if !s.isTableExists(table.TableName) {
			return fmt.Errorf("table '%s' does not exist in the database", table.TableName)
		}

		// Validate the columns in the model
		columns, _ := structToColumnsAndValues(reflect.New(reflect.TypeOf(table.Model)).Interface())
		if err := s.validateColumns(table.TableName, columns); err != nil {
			return fmt.Errorf("validation error for table '%s': %w", table.TableName, err)
		}

		// Start transaction
		tx, err := s.DB.Begin()
		if err != nil {
			return err
		}

		// Seed rows
		for i := 0; i < table.RowCount; i++ {
			// Populate the struct using Faker
			row := reflect.New(reflect.TypeOf(table.Model)).Interface()
			if err := faker.FakeData(row); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to generate fake data: %v", err)
			}

			// Convert struct to columns and value
			columns, values := structToColumnsAndValues(row)

			// Build and execute query
			query := buildInsertQuery(table.TableName, columns)
			if _, err := tx.Exec(query, values...); err != nil {
				tx.Rollback()
				return fmt.Errorf("error inserting data into table '%s': %w", table.TableName, err)
			}
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction for table '%s': %w", table.TableName, err)
		}
	}

	slog.Info("Seeding completed successfully")
	return nil
}

// Helper: Checks if a table exists in the database.
func (s *Seeder) isTableExists(tableName string) bool {
	query := `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)`
	var exists bool
	if err := s.DB.QueryRow(query, tableName).Scan(&exists); err != nil {
		log.Printf("Failed to check table existence: %v\n", err)
		return false
	}
	return exists
}

// Helper: Validates that the columns in the struct exist in the database table.
func (s *Seeder) validateColumns(tableName string, columns []string) error {
	query := `SELECT column_name FROM information_schema.columns WHERE table_name = $1`
	rows, err := s.DB.Query(query, tableName)
	if err != nil {
		return fmt.Errorf("failed to query columns for table '%s': %w", tableName, err)
	}
	defer rows.Close()

	// Collect existing columns
	existingColumns := map[string]struct{}{}
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			return fmt.Errorf("failed to scan column name: %w", err)
		}
		existingColumns[column] = struct{}{}
	}

	// Check for missing columns
	for _, col := range columns {
		if _, exists := existingColumns[col]; !exists {
			return fmt.Errorf("missing column '%s' in table '%s'", col, tableName)
		}
	}

	return nil
}

// Helper: Converts struct to columns and values
func structToColumnsAndValues(model interface{}) ([]string, []interface{}) {
	val := reflect.ValueOf(model).Elem()
	typ := val.Type()

	var columns []string
	var values []interface{}
	for i := 0; i < val.NumField(); i++ {
		columns = append(columns, typ.Field(i).Tag.Get("db"))
		values = append(values, val.Field(i).Interface())
	}

	return columns, values
}

// Helper: Builds SQL Insert query
func buildInsertQuery(tableName string, columns []string) string {
	colNames := "(" + strings.Join(columns, ", ") + ")"
	valPlaceholders := make([]string, len(columns))
	for i := range columns {
		valPlaceholders[i] = fmt.Sprintf("$%d", i+1)
	}
	valPlaceholdersStr := "(" + strings.Join(valPlaceholders, ", ") + ")"
	return fmt.Sprintf("INSERT INTO %s %s VALUES %s", tableName, colNames, valPlaceholdersStr)
}
