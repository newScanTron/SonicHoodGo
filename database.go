package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

type Database struct {
	Name     string
	Db       *sql.DB
	Server   string
	Port     int
	User     string
	Password string
	Database string
}

var dbStct *Database

func (d *Database) getBands() error {
	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	var err error
	bands := make([]string, 0)
	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool:", err.Error())
	}

	ctx := context.Background()

	// Check if database is alive.
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal("Error pinging database: " + err.Error())
	}

	tsql := fmt.Sprintf("SELECT BandID, BandName FROM Band;")

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Error reading rows: " + err.Error())
	}

	defer rows.Close()

	var count int = 0
	// Iterate through the result set.
	for rows.Next() {
		var name string
		var id int

		// Get values from row.
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal("Error reading rows: " + err.Error())
		}

		fmt.Printf("ID: %d, Name: %s\n", id, name)
		count++
	}

	fmt.Printf("count: %d\n", count)
	return err
}
