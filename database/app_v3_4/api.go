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
	sourceUsers any,
	targetUsers any,
) []types.UserIDChange {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.InsertUsersData(db.DB, sourceUsers, targetUsers)
	case types.Postgres:
		return postgres_db.InsertUsersData(db.DB, sourceUsers, targetUsers)
	}

	return nil
}

func InsertGroups(
	db *types.DBConn,
	groups any,
) error {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.InsertGroups(db.DB, groups)
	case types.Postgres:
		return postgres_db.InsertGroups(db.DB, groups)
	}

	return nil
}

func InsertPermissions(
	db *types.DBConn,
	perms any,
) error {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.InsertPermissions(db.DB, perms)
	case types.Postgres:
		return postgres_db.InsertPermissions(db.DB, perms)
	}

	return nil
}

func InsertGroupsPermissions(
	db *types.DBConn,
	gp any,
) error {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.InsertGroupsPermissions(db.DB, gp)
	case types.Postgres:
		return postgres_db.InsertGroupsPermissions(db.DB, gp)
	}

	return nil
}

func InsertDocumentTypes(
	db *types.DBConn,
	v any,
) error {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.InsertDocumentTypes(db.DB, v)
	case types.Postgres:
		return postgres_db.InsertDocumentTypes(db.DB, v)
	}

	return nil
}

func InsertCustomFields(
	db *types.DBConn,
	v any,
) error {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.InsertCustomFields(db.DB, v)
	case types.Postgres:
		return postgres_db.InsertCustomFields(db.DB, v)
	}

	return nil
}

func InsertCustomFieldValues(
	db *types.DBConn,
	v any,
) error {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.InsertCustomFieldValues(db.DB, v)
	case types.Postgres:
		return postgres_db.InsertCustomFieldValues(db.DB, v)
	}

	return nil
}

func InsertDocumentTypesCustomFields(
	db *types.DBConn,
	v any,
) error {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.InsertDocumentTypesCustomFields(db.DB, v)
	case types.Postgres:
		return postgres_db.InsertDocumentTypesCustomFields(db.DB, v)
	}

	return nil
}

func InsertTags(
	db *types.DBConn,
	v any,
) error {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.InsertTags(db.DB, v)
	case types.Postgres:
		return postgres_db.InsertTags(db.DB, v)
	}

	return nil
}

func InsertNodesTags(
	db *types.DBConn,
	v any,
) error {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.InsertNodesTags(db.DB, v)
	case types.Postgres:
		return postgres_db.InsertNodesTags(db.DB, v)
	}

	return nil
}

func InsertUsersGroups(
	db *types.DBConn,
	v any,
) error {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.InsertUsersGroups(db.DB, v)
	case types.Postgres:
		return postgres_db.InsertUsersGroups(db.DB, v)
	}

	return nil
}

func InsertUsersPermissions(
	db *types.DBConn,
	v any,
) error {
	switch db.DBType {
	case types.SQLite:
		return sqlite_db.InsertUsersPermissions(db.DB, v)
	case types.Postgres:
		return postgres_db.InsertUsersPermissions(db.DB, v)
	}

	return nil
}
