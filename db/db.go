package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/lib/pq"
)

type DB struct {
	name     string
	username string
	password string
	sslmode  string
	dialect  string
	conn     *sql.DB
}

func NewDB(name, username string) *DB {
	db := &DB{
		name:     name,
		username: username,
		sslmode:  "disable",
		dialect:  "postgres",
	}
	if err := db.connect(); err != nil {
		fmt.Println("Error connecting to database.")
	}

	return db
}

func (db *DB) connect() error {
	connStr := db.getConnectionString()

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	db.conn = conn

	return nil
}

func (db *DB) getConnectionString() string {
	return fmt.Sprintf("dbname=%s user=%s sslmode=%s",
		db.name,
		db.username,
		db.sslmode,
	)
}

func (db *DB) ExecuteFromFile(filepath string) (string, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	query := string(file)

	result, _ := db.exectute(query)
	return result, nil
}

func (db *DB) exectute(query string) (string, error) {
	results, err := db.conn.Query(query)

	fmt.Println(results.Columns())
	fmt.Println(err)

	return "", err
}
