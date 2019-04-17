package db

import (
	"fmt"
	"github.com/jackc/pgx"
	"os"
)

var DatabaseHost = "192.168.88.26"

//var DatabaseHost = "localhost"
var DatabaseUser = "postgres"
var DatabasePassword = "toor"
var DatabaseDatabase = "pirates"

func Connect() (conn *pgx.Conn, err error) {
	conn, err = pgx.Connect(pgx.ConnConfig{
		Host:     DatabaseHost,
		Port:     5432,
		User:     DatabaseUser,
		Password: DatabasePassword,
		Database: DatabaseDatabase,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}

	return

}
