package main

import (
	"database/sql"
	"log"

	"github.com/ezrantn/goseed"
	_ "github.com/lib/pq"
)

// User defines the structure of the users table.
type user struct {
	ID    string `faker:"uuid_hyphenated" db:"id"`
	Name  string `faker:"name" db:"name"`
	Email string `faker:"email" db:"email"`
}

// Product defines the structure of the products table.
type product struct {
	ID    string  `faker:"uuid_hyphenated" db:"id"`
	Name  string  `faker:"word" db:"name"`
	Price float64 `faker:"amount" db:"price"`
}

// seedConfig defines all table seeders
func seedConfig() []goseed.TableSeeder {
	return []goseed.TableSeeder{
		{
			TableName: "users",
			Model:     user{},
			RowCount:  100,
			BatchSize: 100,
		},
		{
			TableName: "products",
			Model:     product{},
			RowCount:  100,
			BatchSize: 100,
		},
	}
}

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:goseed@localhost:5432/goseed?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	adapter := &goseed.PostgresAdapter{DB: db}

	// Create a new Goseed instance
	seed, err := goseed.NewGoSeed(adapter)
	if err != nil {
		return
	}

	// Add table seeders
	for _, tableSeeders := range seedConfig() {
		seed.Add(tableSeeders)
	}

	// Run the seeder
	if err := seed.Run(); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}
}
