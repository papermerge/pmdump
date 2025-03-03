package database_app_v3_3

import (
	_ "github.com/mattn/go-sqlite3"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/types"
)

func GetGroups(db *types.DBConn) ([]models.Group, error) {
	rows, err := db.DB.Query("SELECT id, name FROM groups")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.Group

	for rows.Next() {
		var group models.Group

		err = rows.Scan(
			&group.ID,
			&group.Name,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}
