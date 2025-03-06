package postgres_app_v3_3

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func Open(dburl string) (*sql.DB, error) {
	return sql.Open("postgres", dburl)
}
