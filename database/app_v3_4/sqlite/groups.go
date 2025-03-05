package sqlite_app_v3_4

import (
	"database/sql"
	"fmt"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/utils"
)

func InsertGroups(db *sql.DB, gr any) error {
	groups := gr.([]models.Group)
	query := `INSERT INTO groups (id, name) VALUES(?, ?)`

	for _, group := range groups {
		uid := utils.UUID2STR(group.ID)
		_, err := db.Exec(query,
			uid,
			group.Name,
		)

		if err != nil {
			return fmt.Errorf(
				"insert group %v failed: %v",
				group,
				err,
			)
		}
	}

	return nil
}
