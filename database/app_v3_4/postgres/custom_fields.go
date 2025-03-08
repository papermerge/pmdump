package postgres_app_v3_4

import (
	"database/sql"
	"fmt"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/utils"
)

func InsertCustomFields(db *sql.DB, cf any) error {
	custom_fields := cf.([]models.CustomField)
	query := `INSERT INTO custom_fields (id, name, type, extra_data, created_at, user_id) VALUES($1, $2, $3, $4, $5, $6)`

	for _, custom_field := range custom_fields {
		id := utils.UUID2STR(custom_field.ID)
		user_id := utils.UUID2STR(custom_field.UserID)
		_, err := db.Exec(query,
			id,
			custom_field.Name,
			custom_field.Type,
			custom_field.ExtraData,
			custom_field.CreatedAt,
			user_id,
		)

		if err != nil {
			return fmt.Errorf(
				"insert custom_field %v failed: %v",
				custom_field,
				err,
			)
		}
	}

	return nil
}
