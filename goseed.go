package goseed

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/fatih/color"
	"github.com/go-faker/faker/v4"
)

/*
	This file contains an implementation of goseed, a database seeder library for Go.
	Developed by: @ezrantn
	GitHub: github.com/ezrantn/goseed
	Contact: ezrantn@proton.me

	This library is licensed under the BSD-3 Clause. See LICENSE for more details.
*/

type DBAdapter interface {
	Ping() error
	IsTableExists(tableName string) (bool, error)
	InsertRow(tableName string, columns []string, values [][]any) error
	GetColumns(tableName string) ([]string, error)
}

// Defines the structure for seeding a specific table
type TableSeeder struct {
	TableName string
	RowCount  int
	Model     any
	BatchSize int
}

// Seeder orchestrates the database seeding process
type Seeder struct {
	Adapter      DBAdapter
	TableSeeders []TableSeeder
}

// Creates a new seeder instance
func NewGoSeed(adapter DBAdapter) (*Seeder, error) {
	if adapter == nil {
		return nil, errors.New("database connection is nil")
	}

	if err := adapter.Ping(); err != nil {
		return nil, fmt.Errorf("goseed unable to connect to database: %v", err)
	}

	return &Seeder{
		Adapter:      adapter,
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
		fmt.Printf("Seeding %d rows for table '%s'...\n", table.RowCount, table.TableName)

		// Check if the table exists
		exists, err := s.Adapter.IsTableExists(table.TableName)
		if err != nil || !exists {
			s.logError(fmt.Sprintf("Table '%s' does not exist in the database", table.TableName))
			return fmt.Errorf("table '%s' does not exist in the database", table.TableName)
		}

		// Validate the columns in the model
		columns, err := s.Adapter.GetColumns(table.TableName)
		if err != nil {
			s.logError(fmt.Sprintf("Column validation failed for table '%s': %v", table.TableName, err))
			return err
		}

		// Batch size for inserting
		batchSize := table.BatchSize
		var valuesBatch [][]any

		// Seed rows
		for i := 0; i < table.RowCount; i++ {
			// Populate the struct using Faker
			row := reflect.New(reflect.TypeOf(table.Model)).Interface()
			if err := faker.FakeData(&row); err != nil {
				return fmt.Errorf("failed to generate fake data for %s: %v", table.TableName, err)
			}

			// Convert struct to columns and value
			_, values := structToColumnsAndValues(row)
			valuesBatch = append(valuesBatch, values)

			// Insert in batches
			if len(valuesBatch) >= batchSize || i == table.RowCount-1 {
				// Use batch insert method specific to the adapter
				if postgresAdapter, ok := s.Adapter.(*PostgresAdapter); ok {
					if err := postgresAdapter.InsertRow(table.TableName, columns, valuesBatch); err != nil {
						s.logError(fmt.Sprintf("Error inserting batch into table '%s': %v ", table.TableName, err))
						return err
					}
				} else if mysqlAdapter, ok := s.Adapter.(*MySQLAdapter); ok {
					if err := mysqlAdapter.InsertRow(table.TableName, columns, valuesBatch); err != nil {
						s.logError(fmt.Sprintf("Error inserting batch into table '%s': %v", table.TableName, err))
						return err
					}
				}

				// Reset batch
				valuesBatch = [][]any{}
			}
		}
	}

	fmt.Println("Seeding completed successfully")
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

// Helper: Print out error log
func (s *Seeder) logError(message string) {
	red := color.RedString("Error:")
	fmt.Printf("%s %s\n", red, message)
}
