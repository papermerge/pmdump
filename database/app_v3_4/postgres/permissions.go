package postgres_app_v3_4

import (
	"database/sql"
	"fmt"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/utils"
)

func InsertPermissions(db *sql.DB, p any) error {
	permissions := p.([]models.Permission)
	query := `INSERT INTO permissions (id, name, codename) VALUES($1, $2, $3)`

	for _, perm := range permissions {
		uid := utils.UUID2STR(perm.ID)
		_, err := db.Exec(query,
			uid,
			perm.Name,
			perm.Codename,
		)

		if err != nil {
			return fmt.Errorf(
				"insert permission %v failed: %v",
				perm,
				err,
			)
		}
	}

	return nil
}
