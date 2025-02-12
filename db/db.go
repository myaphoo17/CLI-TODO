package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
)

var Db *sql.DB

func InitDB(dsn string) {
	var err error
	Db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		fmt.Println("Please run the main program to create the database.")
		os.Exit(1)
		CreateAndInitializeDatabase()
	}

	if err = Db.Ping(); err != nil {
		fmt.Println("Database not found, attempting to create it...")

		tempDb, err := sql.Open("mysql", dsn[:strings.LastIndex(dsn, "/")+1])
		if err != nil {
			fmt.Println("Error connecting to the MySQL server:", err)
			os.Exit(1)
		}
		defer tempDb.Close()

		CreateDatabase(tempDb)

		cfg, err := ini.Load("config.ini")
		if err != nil {
			fmt.Println("Failed to read config file:", err)
			os.Exit(1)
		}

		dbName := cfg.Section("database").Key("dbName").String()
		if dbName == "" {
			fmt.Println("Database name is empty in config file")
			os.Exit(1)
		}
		dsn = dsn[:strings.LastIndex(dsn, "/")+1] + dbName

		Db, err = sql.Open("mysql", dsn)
		if err != nil {
			fmt.Println("Error reconnecting to the newly created database:", err)
			os.Exit(1)
		}
		if err = Db.Ping(); err != nil {
			fmt.Println("Failed to connect to the database after creation:", err)
			os.Exit(1)
		}

		CreateAndInitializeDatabase()
	}
}

func CreateDatabase(tempDb *sql.DB) {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Println("Failed to read config file:", err)
		os.Exit(1)
	}

	dbName := cfg.Section("database").Key("dbName").String()
	if dbName == "" {
		fmt.Println("Database name is empty in config file")
		os.Exit(1)
	}

	_, err = tempDb.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
	if err != nil {
		fmt.Println("Error creating database:", err)
		os.Exit(1)
	}
	fmt.Println("Database created successfully!")
}

func CreateAndInitializeDatabase() {
	_, err := Db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		fmt.Println("Error creating tasks table:", err)
		fmt.Println("Please check your database configuration and try again.")
		os.Exit(1)
	}
	fmt.Println("Database initialized successfully!")
}
