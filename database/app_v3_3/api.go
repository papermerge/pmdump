package database_app_v3_3

import (
	"fmt"
	"net/url"
	"strings"

	postgres_db "github.com/papermerge/pmdump/database/app_v3_3/postgres"
	sqlite_db "github.com/papermerge/pmdump/database/app_v3_3/sqlite"

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

	return nil, fmt.Errorf("database open: app version %q not supported", appVer)
}

func GetUsers(db *types.DBConn) (interface{}, error) {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.GetUsers(db.DB)
	case types.Postgres:
		return postgres_db.GetUsers(db.DB)
	}

	return nil, fmt.Errorf("database GetUsers: db type %q not supported", db.DBType)
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

func GetDocumentPageRows(db *types.DBConn, user_id interface{}) (interface{}, error) {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.GetDocumentPageRows(db.DB, user_id)
	case types.Postgres:
		return postgres_db.GetDocumentPageRows(db.DB, user_id)
	}

	err := fmt.Errorf(
		"database GetDocumentPageRows: db type %q not supported",
		db.DBType,
	)

	return nil, err
}
