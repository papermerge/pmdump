package sqlite_app_v3_4

import (
	"database/sql"
	"fmt"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/utils"
)

func InsertUsersPermissions(db *sql.DB, p any) error {
	users_permissions := p.([]models.UsersPermissions)
	query := `INSERT INTO users_permissions (user_id, permission_id) VALUES(?, ?)`

	for _, user_perm := range users_permissions {
		uid := utils.UUID2STR(user_perm.UserID)
		puid := utils.UUID2STR(user_perm.PermissionID)
		_, err := db.Exec(query,
			uid,
			puid,
		)

		if err != nil {
			return fmt.Errorf(
				"insert user permission %v failed: %v",
				user_perm,
				err,
			)
		}
	}

	return nil
}
