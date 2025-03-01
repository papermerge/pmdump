package sqlite_app_v3_3

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	models "github.com/papermerge/pmdump/models/app_v3_3"
)

func GetUsers(db *sql.DB) (models.Users, error) {
	rows, err := db.Query("SELECT id, username, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.UUID, &user.Username, &user.EMail)
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

func GetTargetUsers(db *sql.DB) (models.TargetUserList, error) {
	rows, err := db.Query("SELECT id, username, email, home_folder_id, inbox_folder_id FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users models.TargetUserList

	for rows.Next() {
		var user models.TargetUser
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.EMail,
			&user.HomeID,
			&user.InboxID,
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

func InsertUsersData(
	db *sql.DB,
	sourceUsers []models.User,
	targetUsers models.TargetUserList,
) {

	for i := 0; i < len(sourceUsers); i++ {
		targetUser := targetUsers.Get(sourceUsers[i].Username)
		if targetUser != nil {
			ImportUserData(db, sourceUsers[i], targetUser)
		} else {
			targetUser, err := CreateTargetUser(db, sourceUsers[i])
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"Error creating target user for %s: %v\n",
					sourceUsers[i].Username,
					err,
				)
				continue
			}
			ImportUserData(db, sourceUsers[i], targetUser)
		}
	}
}

func ImportUserData(
	db *sql.DB,
	sourceUser models.User,
	targetUser *models.TargetUser,
) {

	models.ForEachNode(sourceUser.Home, models.UpdateNodeUUID)
	models.ForEachNode(sourceUser.Inbox, models.UpdateNodeUUID)

	ForEachSourceNode(
		db,
		sourceUser.Home, // start here
		targetUser.HomeID,
		targetUser.ID,
		CreateTargetNode,
	)
	ForEachSourceNode(
		db,
		sourceUser.Inbox, // start here
		targetUser.InboxID,
		targetUser.ID,
		CreateTargetNode,
	)
}

func CreateTargetUser(
	db *sql.DB,
	source models.User,
) (*models.TargetUser, error) {
	return nil, nil
}
