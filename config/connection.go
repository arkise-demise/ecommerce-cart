package config

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var dbInstance *sql.DB
var once sync.Once

func Connection() (*sql.DB, error) {
	var err error

	once.Do(func() {
		cfg, err := LoadConfig()
		if err != nil {
			err = fmt.Errorf("error loading config: %v", err)
			return
		}

		dsn := cfg.DSN()
		dbInstance, err = sql.Open("postgres", dsn)
		if err != nil {
			err = fmt.Errorf("error opening database connection: %v", err)
			return
		}

		dbInstance.SetMaxOpenConns(20)
		dbInstance.SetMaxIdleConns(10)
		dbInstance.SetConnMaxLifetime(5 * time.Minute)

		if err = dbInstance.Ping(); err != nil {
			err = fmt.Errorf("error pinging database: %v", err)
			return
		}

		log.Println("Database connection established successfully.")
	})

	if err != nil {
		return nil, err
	}

	return dbInstance, nil
}
