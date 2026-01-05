package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Convert mysql:// URL to DSN format
	// mysql://user:pass@host:port/dbname -> user:pass@tcp(host:port)/dbname
	db, err := sql.Open("mysql", parseDSN(dsn))
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected")
	return db
}

func parseDSN(url string) string {
	// mysql://user:pass@host:port/dbname -> user:pass@tcp(host:port)/dbname?parseTime=true
	// Simple parser for mysql:// URLs
	if len(url) > 8 && url[:8] == "mysql://" {
		url = url[8:]
	}

	// Find @ to split credentials and host
	atIdx := -1
	for i, c := range url {
		if c == '@' {
			atIdx = i
			break
		}
	}

	if atIdx == -1 {
		return url + "?parseTime=true"
	}

	creds := url[:atIdx]
	rest := url[atIdx+1:]

	// Find / to split host and dbname
	slashIdx := -1
	for i, c := range rest {
		if c == '/' {
			slashIdx = i
			break
		}
	}

	if slashIdx == -1 {
		return fmt.Sprintf("%s@tcp(%s)/?parseTime=true", creds, rest)
	}

	host := rest[:slashIdx]
	dbname := rest[slashIdx+1:]

	return fmt.Sprintf("%s@tcp(%s)/%s?parseTime=true", creds, host, dbname)
}
