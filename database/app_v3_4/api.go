package database_app_v3_3

import (
	"fmt"
	"net/url"
	"strings"

	postgres_db "github.com/papermerge/pmdump/database/app_v3_4/postgres"
	sqlite_db "github.com/papermerge/pmdump/database/app_v3_4/sqlite"
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

func GetTargetUsers(db *types.DBConn) (interface{}, error) {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.GetTargetUsers(db.DB)
	case types.Postgres:
		return postgres_db.GetTargetUsers(db.DB)
	}

	return nil, fmt.Errorf("database GetTargetUsers: db type %q not supported", db.DBType)
}

func InsertUsersData(
	db *types.DBConn,
	sourceUsers interface{},
	targetUsers interface{},
) {
	switch db.DBType {
	case types.SQLite:
		sqlite_db.InsertUsersData(db.DB, sourceUsers, targetUsers)
	case types.Postgres:
		postgres_db.InsertUsersData(db.DB, sourceUsers, targetUsers)
	}
}
