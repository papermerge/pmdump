package postgres_app_v3_4

import (
	"database/sql"
	"fmt"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/utils"
)

func InsertTags(db *sql.DB, t any) error {
	tags := t.([]models.Tag)
	query := `
    INSERT INTO tags (
      id,
      name,
      fg_color,
      bg_color,
      pinned,
      description,
      user_id
    ) VALUES(
      ?, ?, ?,
      ?, ?, ?,
      ?
    )`

	for _, tag := range tags {
		id := utils.UUID2STR(tag.ID)
		user_id := utils.UUID2STR(tag.UserID)
		_, err := db.Exec(query,
			id,
			tag.Name,
			tag.FGColor,
			tag.BGColor,
			tag.Pinned,
			tag.Description,
			user_id,
		)

		if err != nil {
			return fmt.Errorf(
				"insert tag %v failed: %v",
				tag,
				err,
			)
		}
	}

	return nil
}
