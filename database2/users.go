package database2

import (
	"database/sql"

	"github.com/papermerge/pmg-dump/models2"
)

func GetUsers(db *sql.DB) ([]models2.User, error) {
	rows, err := db.Query("SELECT id, username, email FROM core_user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models2.User

	for rows.Next() {
		var user models2.User
		err = rows.Scan(&user.ID, &user.Username, &user.EMail)
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
