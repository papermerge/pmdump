package database_app_v3_3

import (
	_ "github.com/mattn/go-sqlite3"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/types"
)

func GetTags(db *types.DBConn) ([]models.Tag, error) {
	rows, err := db.DB.Query("SELECT id, name, fg_color, bg_color, pinned, description, user_id FROM tags")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.Tag

	for rows.Next() {
		var entry models.Tag

		err = rows.Scan(
			&entry.ID,
			&entry.Name,
			&entry.FGColor,
			&entry.BGColor,
			&entry.Pinned,
			&entry.Description,
			&entry.UserID,
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
