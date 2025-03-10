package database_app_v2_1

import (
	"fmt"
	"net/url"
	"strings"

	postgres_db "github.com/papermerge/pmdump/database/app_v3_2/postgres"
	sqlite_db "github.com/papermerge/pmdump/database/app_v3_2/sqlite"

	"github.com/papermerge/pmdump/types"
)

func Open(dburl string, appVer types.AppVersion) (*types.DBConn, error) {
	parsedDBURL, err := url.Parse(dburl)

	if err != nil {
		return nil, fmt.Errorf("Error parsing dburl %s: %v", dburl, err)
	}

	if strings.HasPrefix(parsedDBURL.Scheme, "sqlite") {
		db, err := sqlite_db.Open(parsedDBURL.Path)
		if err != nil {
			return nil, err
		}
		dbconn := types.DBConn{
			AppVersion: appVer,
			DBType:     types.SQLite,
			DB:         db,
		}

		return &dbconn, nil
	}

	if strings.HasPrefix(parsedDBURL.Scheme, "postgres") {
		db, err := postgres_db.Open(dburl)
		if err != nil {
			return nil, err
		}
		dbconn := types.DBConn{
			AppVersion: appVer,
			DBType:     types.Postgres,
			DB:         db,
		}

		return &dbconn, nil
	}

	return nil, fmt.Errorf("database open: app version %q not supported", appVer)
}

func GetHomeFlatNodes(db *types.DBConn, user_id interface{}) (interface{}, error) {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.GetHomeFlatNodes(db.DB, user_id)
	case types.Postgres:
		return postgres_db.GetHomeFlatNodes(db.DB, user_id)
	}

	return nil, fmt.Errorf("database GetHomeFlatNodes: db type %q not supported", db.DBType)
}

func GetInboxFlatNodes(db *types.DBConn, user_id interface{}) (interface{}, error) {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.GetInboxFlatNodes(db.DB, user_id)
	case types.Postgres:
		return postgres_db.GetInboxFlatNodes(db.DB, user_id)
	}

	err := fmt.Errorf(
		"database GetInboxFlatNodes: db type %q not supported",
		db.DBType,
	)

	return nil, err
}

func GetUserNodes(db *types.DBConn, user *interface{}) error {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.GetUserNodes(db.DB, user)
	case types.Postgres:
		return postgres_db.GetUserNodes(db.DB, user)
	}

	return fmt.Errorf(
		"database GetUserNodes: db type %q not supported",
		db.DBType,
	)
}

func InsertDocVersionsAndPages(db *types.DBConn, node any) error {
	switch db.DBType {
	case types.SQLite:
		sqlite_db.InsertDocVersionsAndPages(db.DB, node)
	case types.Postgres:
		postgres_db.InsertDocVersionsAndPages(db.DB, node)
	}

	return fmt.Errorf(
		"database GetDocumentPageRows: db type %q not supported",
		db.DBType,
	)

}
