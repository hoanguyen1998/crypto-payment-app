package drivers

import (
	"database/sql"
	"log"
	"time"
)

func ConnectPostgres(connectStr string) (*sql.DB, error) {
	// db, err := sql.Open("sqlite3", "crypto.db")
	db, err := sql.Open("pgx", connectStr)

	if err != nil {
		log.Println("Error occurred")
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
