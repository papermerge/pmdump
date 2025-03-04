package database_app_v3_3

import (
	_ "github.com/mattn/go-sqlite3"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/types"
)

func GetUsersPermissions(db *types.DBConn) ([]models.UsersPermissions, error) {
	query := `
    SELECT
      user_id,
      permission_id
    FROM users_permissions
  `
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.UsersPermissions

	for rows.Next() {
		var entry models.UsersPermissions

		err = rows.Scan(
			&entry.UserID,
			&entry.PermissionID,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
