package database_app_v3_3

import (
	_ "github.com/mattn/go-sqlite3"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/types"
)

func GetGroupsPermissions(db *types.DBConn) ([]models.GroupsPermissions, error) {
	rows, err := db.DB.Query("SELECT group_id, permission_id FROM groups_permissions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.GroupsPermissions

	for rows.Next() {
		var entry models.GroupsPermissions

		err = rows.Scan(
			&entry.GroupID,
			&entry.PermissionID,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
