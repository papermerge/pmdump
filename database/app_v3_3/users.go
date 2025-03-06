package database_app_v3_3

import (
	"database/sql"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	models "github.com/papermerge/pmdump/models/app_v3_3"
)

func GetUsers(db *sql.DB) ([]models.User, error) {
	rows, err := db.Query("SELECT id, home_folder_id, inbox_folder_id, username, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err = rows.Scan(
			&user.ID,
			&user.HomeFolderID,
			&user.InboxFolderID,
			&user.Username,
			&user.EMail,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
