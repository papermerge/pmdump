package sqlite_app_v3_4

import (
	"database/sql"
	"fmt"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/utils"
)

func InsertUsersGroups(db *sql.DB, p any) error {
	users_groups := p.([]models.UsersGroups)
	query := `INSERT INTO users_groups (user_id, group_id) VALUES(?, ?)`

	for _, user_group := range users_groups {
		uid := utils.UUID2STR(user_group.UserID)
		guid := utils.UUID2STR(user_group.GroupID)
		_, err := db.Exec(query,
			uid,
			guid,
		)

		if err != nil {
			return fmt.Errorf(
				"insert user group %v failed: %v",
				user_group,
				err,
			)
		}
	}

	return nil
}
