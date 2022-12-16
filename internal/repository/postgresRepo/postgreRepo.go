package postgresRepo

import (
	"database/sql"

	"github.com/hoanguyen1998/crypto-payment-system/internal/repository"
	_ "github.com/mattn/go-sqlite3"
)

type postgresDBRepo struct {
	DB *sql.DB
}

func NewPostgresRepo(db *sql.DB) repository.DatabaseRepo {
	return &postgresDBRepo{
		DB: db,
	}
}
