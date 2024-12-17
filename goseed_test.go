package goseed

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type MockDBAdapter struct {
	PingFunc        func() error
	IsTableExistsFn func(tableName string) (bool, error)
	InsertRowFn     func(tableName string, columns []string, values []any) error
	GetColumnsFn    func(tableName string) ([]string, error)
}

func (m *MockDBAdapter) Ping() error {
	if m.PingFunc != nil {
		return m.PingFunc()
	}

	return nil
}

func (m *MockDBAdapter) IsTableExists(tableName string) (bool, error) {
	if m.IsTableExistsFn != nil {
		return m.IsTableExistsFn(tableName)
	}

	return true, nil
}

func (m *MockDBAdapter) InsertRow(tableName string, columns []string, values []any) error {
	if m.InsertRowFn != nil {
		return m.InsertRowFn(tableName, columns, values)
	}

	return nil
}

func (m *MockDBAdapter) GetColumns(tableName string) ([]string, error) {
	if m.GetColumnsFn != nil {
		return m.GetColumnsFn(tableName)
	}

	return []string{"id", "name", "age"}, nil
}

// Test struct to be used as a model
type TestUser struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

// TestSeeder validates the seeding process
func TestSeeder(t *testing.T) {
	mockDB := &MockDBAdapter{
		PingFunc: func() error {
			return nil
		},
		IsTableExistsFn: func(tableName string) (bool, error) {
			if tableName == "users" {
				return true, nil
			}
			return false, fmt.Errorf("table %s does not exist", tableName)
		},
		InsertRowFn: func(tableName string, columns []string, values []any) error {
			if tableName != "users" {
				return errors.New("invalid table")
			}

			if !reflect.DeepEqual(columns, []string{"id", "name", "age"}) {
				return errors.New("column mismatch")
			}

			return nil
		},
		GetColumnsFn: func(tableName string) ([]string, error) {
			return []string{"id", "name", "age"}, nil
		},
	}

	// Step 1: Create a new seeder
	seeder, err := NewGoSeed(mockDB)
	if err != nil {
		t.Fatalf("failed to create seeder: %v", err)
	}

	// Step 2: Add a TableSeeder
	err = seeder.Add(TableSeeder{
		TableName: "users",
		RowCount:  3,
		Model:     TestUser{},
	})

	if err != nil {
		t.Fatalf("failed to add table seeder: %v", err)
	}

	// Step 3: Run the seeder
	err = seeder.Run()
	if err != nil {
		t.Fatalf("seeder failed to run: %v", err)
	}

	t.Log("Seeder ran successfully")
}

// TestInvalidSeeder verifies handling of invalid inputs
func TestInvalidSeeder(t *testing.T) {
	mockDB := &MockDBAdapter{
		PingFunc: func() error {
			return nil
		},
	}

	// Step 1: Create a seeder and test invalid TableSeeder
	_, err := NewGoSeed(nil)
	if err == nil {
		t.Fatalf("expected error when adapter is nil, get nil")
	}

	// Step 2: Create a seeder and test invalid TableSeeder
	seeder, err := NewGoSeed(mockDB)
	if err != nil {
		t.Fatalf("failed to create seeder: %v", err)
	}

	invalidSeeders := []TableSeeder{
		{TableName: "", RowCount: 3, Model: TestUser{}},
		{TableName: "users", RowCount: 0, Model: TestUser{}},
		{TableName: "users", RowCount: 3, Model: ""},
	}

	for _, seederConfig := range invalidSeeders {
		err := seeder.Add(seederConfig)
		if err == nil {
			t.Fatalf("expected error for invalid seeder: %+v, got nil", seederConfig)
		}
	}

	t.Log("invalid seeder inputs handled successfully")
}
