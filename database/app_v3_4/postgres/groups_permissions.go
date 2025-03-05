package postgres_app_v3_4

import (
	"database/sql"
	"fmt"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/utils"
)

func InsertGroupsPermissions(db *sql.DB, p any) error {
	groups_permissions := p.([]models.GroupsPermissions)
	query := `INSERT INTO groups_permissions (group_id, permission_id) VALUES(?, ?)`

	for _, group_perm := range groups_permissions {
		guid := utils.UUID2STR(group_perm.GroupID)
		puid := utils.UUID2STR(group_perm.PermissionID)
		_, err := db.Exec(query,
			guid,
			puid,
		)

		if err != nil {
			return fmt.Errorf(
				"insert group permission %v failed: %v",
				group_perm,
				err,
			)
		}
	}

	return nil
}
