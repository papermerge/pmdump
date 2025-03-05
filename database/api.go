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
	case types.V3_3, types.V3_4:
		return database_app_v3_3.Open(dburl, appVer)
	}

	return nil, fmt.Errorf("database open: app version %q not supported", appVer)
}

func GetUsers(db *types.DBConn) (any, error) {
	switch db.AppVersion {
	case types.V2_0:
		return database_app_v2_0.GetUsers(db)
	case types.V3_3:
		return database_app_v3_3.GetUsers(db)
	}

	return nil, fmt.Errorf("database GetUsers: app version %q not supported", db.AppVersion)
}

func GetGroups(db *types.DBConn) (any, error) {
	switch db.AppVersion {
	case types.V3_3:
		return database_app_v3_3.GetGroups(db)
	}

	return nil, fmt.Errorf("database GetGroups: app version %q not supported", db.AppVersion)
}

func GetPermissions(db *types.DBConn) (any, error) {
	switch db.AppVersion {
	case types.V3_3:
		return database_app_v3_3.GetPermissions(db)
	}

	return nil, fmt.Errorf("database GetGroups: app version %q not supported", db.AppVersion)
}

func GetGroupsPermissions(db *types.DBConn) (any, error) {
	switch db.AppVersion {
	case types.V3_3:
		return database_app_v3_3.GetGroupsPermissions(db)
	}

	return nil, fmt.Errorf("database GetGroupsPermissions: app version %q not supported", db.AppVersion)
}

func GetDocumentTypes(db *types.DBConn) (any, error) {
	switch db.AppVersion {
	case types.V3_3:
		return database_app_v3_3.GetDocumentTypes(db)
	}

	return nil, fmt.Errorf("database GetGroupsPermissions: app version %q not supported", db.AppVersion)
}

func GetTags(db *types.DBConn) (any, error) {
	switch db.AppVersion {
	case types.V3_3:
		return database_app_v3_3.GetTags(db)
	}

	return nil, fmt.Errorf("database GetTags: app version %q not supported", db.AppVersion)
}

func GetNodesTags(db *types.DBConn) (any, error) {
	switch db.AppVersion {
	case types.V3_3:
		return database_app_v3_3.GetNodesTags(db)
	}

	return nil, fmt.Errorf("database GetNodesTags: app version %q not supported", db.AppVersion)
}

func GetUsersGroups(db *types.DBConn) (any, error) {
	switch db.AppVersion {
	case types.V3_3:
		return database_app_v3_3.GetUsersGroups(db)
	}

	return nil, fmt.Errorf("database GetUserGroups: app version %q not supported", db.AppVersion)
}

func GetUsersPermissions(db *types.DBConn) (any, error) {
	switch db.AppVersion {
	case types.V3_3:
		return database_app_v3_3.GetUsersPermissions(db)
	}

	return nil, fmt.Errorf("database GetUsersPermissions: app version %q not supported", db.AppVersion)
}

func GetCustomFields(db *types.DBConn) (any, error) {
	switch db.AppVersion {
	case types.V3_3:
		return database_app_v3_3.GetCustomFields(db)
	}

	return nil, fmt.Errorf("database GetCustomFields: app version %q not supported", db.AppVersion)
}

func GetDocumentTypesCustomFields(db *types.DBConn) (any, error) {
	switch db.AppVersion {
	case types.V3_3:
		return database_app_v3_3.GetDocumentTypesCustomFields(db)
	}

	return nil, fmt.Errorf("database GetDocumentTypesCustomFields: app version %q not supported", db.AppVersion)
}

func GetCustomFieldValues(db *types.DBConn) (any, error) {
	switch db.AppVersion {
	case types.V3_3:
		return database_app_v3_3.GetCustomFieldValues(db)
	}

	return nil, fmt.Errorf("database GetCustomFieldValues: app version %q not supported", db.AppVersion)
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

func InsertUsersData(db *types.DBConn, sourceUsers any, targetUsers any) ([]types.UserIDChange, error) {
	switch db.AppVersion {
	case types.V3_4:
		return database_app_v3_4.InsertUsersData(db, sourceUsers, targetUsers), nil
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

func GetDocumentPageRows(db *types.DBConn, user_id any) (any, error) {
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

func InsertGroups(db *types.DBConn, groups any) error {
	switch db.AppVersion {
	case types.V3_4:
		return database_app_v3_4.InsertGroups(db, groups)
	}

	err := fmt.Errorf(
		"database InsertGroups: app version %q not supported",
		db.AppVersion,
	)

	return err
}

func InsertPermissions(db *types.DBConn, p any) error {
	switch db.AppVersion {
	case types.V3_4:
		return database_app_v3_4.InsertPermissions(db, p)
	}

	err := fmt.Errorf(
		"database InsertPermissions: app version %q not supported",
		db.AppVersion,
	)

	return err
}

func InsertGroupsPermissions(db *types.DBConn, p any) error {
	switch db.AppVersion {
	case types.V3_4:
		return database_app_v3_4.InsertGroupsPermissions(db, p)
	}

	err := fmt.Errorf(
		"database InsertGroupsPermissions: app version %q not supported",
		db.AppVersion,
	)

	return err
}

func InsertDocumentTypes(db *types.DBConn, p any) error {
	switch db.AppVersion {
	case types.V3_4:
		return database_app_v3_4.InsertDocumentTypes(db, p)
	}

	err := fmt.Errorf(
		"database InsertDocumentTypes: app version %q not supported",
		db.AppVersion,
	)

	return err
}

func InsertCustomFields(db *types.DBConn, p any) error {
	switch db.AppVersion {
	case types.V3_4:
		return database_app_v3_4.InsertCustomFields(db, p)
	}

	err := fmt.Errorf(
		"database InsertCustomFields: app version %q not supported",
		db.AppVersion,
	)

	return err
}

func InsertCustomFieldValues(db *types.DBConn, p any) error {
	switch db.AppVersion {
	case types.V3_4:
		return database_app_v3_4.InsertCustomFieldValues(db, p)
	}

	err := fmt.Errorf(
		"database InsertCustomFieldValues: app version %q not supported",
		db.AppVersion,
	)

	return err
}

func InsertDocumentTypesCustomFields(db *types.DBConn, p any) error {
	switch db.AppVersion {
	case types.V3_4:
		return database_app_v3_4.InsertDocumentTypesCustomFields(db, p)
	}

	err := fmt.Errorf(
		"database InsertDocumentTypesCustomFields: app version %q not supported",
		db.AppVersion,
	)

	return err
}

func InsertTags(db *types.DBConn, p any) error {
	switch db.AppVersion {
	case types.V3_4:
		return database_app_v3_4.InsertTags(db, p)
	}

	err := fmt.Errorf(
		"database Tags: app version %q not supported",
		db.AppVersion,
	)

	return err
}

func InsertNodesTags(db *types.DBConn, p any) error {
	switch db.AppVersion {
	case types.V3_4:
		return database_app_v3_4.InsertNodesTags(db, p)
	}

	err := fmt.Errorf(
		"database NodesTags: app version %q not supported",
		db.AppVersion,
	)

	return err
}

func InsertUsersGroups(db *types.DBConn, p any) error {
	switch db.AppVersion {
	case types.V3_4:
		return database_app_v3_4.InsertUsersGroups(db, p)
	}

	err := fmt.Errorf(
		"database UsersGroups: app version %q not supported",
		db.AppVersion,
	)

	return err
}

func InsertUsersPermissions(db *types.DBConn, p any) error {
	switch db.AppVersion {
	case types.V3_4:
		return database_app_v3_4.InsertUsersPermissions(db, p)
	}

	err := fmt.Errorf(
		"database UsersPermissions: app version %q not supported",
		db.AppVersion,
	)

	return err
}
