package database_app_v3_3

import (
	_ "github.com/mattn/go-sqlite3"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/types"
)

func GetUsersGroups(db *types.DBConn) ([]models.UsersGroups, error) {
	query := `
    SELECT
      group_id,
      user_id
    FROM users_groups
  `
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.UsersGroups

	for rows.Next() {
		var entry models.UsersGroups

		err = rows.Scan(
			&entry.GroupID,
			&entry.UserID,
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
