package database

import (
	"fmt"

	database_app_v2_0 "github.com/papermerge/pmdump/database/app_v2_0"
	database_app_v3_3 "github.com/papermerge/pmdump/database/app_v3_3"
	database_app_v3_4 "github.com/papermerge/pmdump/database/app_v3_4"

	"github.com/papermerge/pmdump/types"
)

func Open(dburl string, appVer types.AppVersion) (*types.DBConn, error) {
	switch appVer {
	case types.V2_0:
		return database_app_v2_0.Open(dburl, appVer)
	case types.V3_3:
		return database_app_v3_3.Open(dburl, appVer)
	}

	return nil, fmt.Errorf("database open: app version %q not supported", appVer)
}

func GetUsers(db *types.DBConn) (interface{}, error) {
	switch db.AppVersion {
	case types.V2_0:
		return database_app_v2_0.GetUsers(db)
	case types.V3_3:
		return database_app_v3_3.GetUsers(db)
	}

	return nil, fmt.Errorf("database GetUsers: app version %q not supported", db.AppVersion)
}

func GetHomeFlatNodes(db *types.DBConn, user_id interface{}) (interface{}, error) {
	switch db.AppVersion {
	case types.V2_0:
		return database_app_v2_0.GetHomeFlatNodes(db, user_id)
	case types.V3_3:
		return database_app_v3_3.GetHomeFlatNodes(db, user_id)
	}

	return nil, fmt.Errorf("database GetHomeFlatNodes: app version %q not supported", db.AppVersion)
}

func GetInboxFlatNodes(db *types.DBConn, user_id interface{}) (interface{}, error) {
	switch db.AppVersion {
	case types.V2_0:
		return database_app_v2_0.GetInboxFlatNodes(db, user_id)
	case types.V3_3:
		return database_app_v3_3.GetInboxFlatNodes(db, user_id)
	}

	return nil, fmt.Errorf("database GetInboxFlatNodes: app version %q not supported", db.AppVersion)
}

func GetUserNodes(db *types.DBConn, user interface{}) error {
	switch db.AppVersion {
	case types.V2_0:
		return database_app_v2_0.GetUserNodes(db, &user)
	case types.V3_3:
		return database_app_v3_3.GetUserNodes(db, &user)
	}

	return fmt.Errorf("database GetUserNodes: app version %q not supported", db.AppVersion)
}

func GetTargetUsers(db *types.DBConn) (interface{}, error) {
	switch db.AppVersion {
	case types.V3_4:
		return database_app_v3_4.GetTargetUsers(db)
	}
	return nil, fmt.Errorf("database GetTargetUsers: app version %q not supported", db.AppVersion)
}

func InsertUsersData(db *types.DBConn, sourceUsers interface{}, targetUsers interface{}) (interface{}, error) {
	switch db.AppVersion {
	case types.V3_4:
		database_app_v3_4.InsertUsersData(db, sourceUsers, targetUsers)
	}
	return nil, fmt.Errorf("database InsertUserData: app version %q not supported", db.AppVersion)
}

func InsertDocVersionsAndPages(db *types.DBConn, node any) error {
	switch db.AppVersion {
	case types.V3_3:
		database_app_v3_3.InsertDocVersionsAndPages(db, node)
	}

	return fmt.Errorf("database InsertDocVersionsAndPages: app version %q not supported", db.AppVersion)
}

func GetDocumentPageRows(db *types.DBConn, user_id interface{}) (interface{}, error) {
	switch db.AppVersion {
	case types.V2_0:
		return database_app_v2_0.GetDocumentPageRows(db, user_id)
	}

	err := fmt.Errorf(
		"database GetDocumentPageRows: app version %q not supported",
		db.AppVersion,
	)

	return nil, err
}
