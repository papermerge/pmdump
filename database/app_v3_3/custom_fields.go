package database_app_v3_3

import (
	_ "github.com/mattn/go-sqlite3"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/types"
)

func GetCustomFields(db *types.DBConn) ([]models.CustomField, error) {
	query := `
    SELECT
      id,
      name,
      type,
      extra_data,
      created_at,
      user_id
    FROM custom_fields
  `
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.CustomField

	for rows.Next() {
		var entry models.CustomField

		err = rows.Scan(
			&entry.ID,
			&entry.Name,
			&entry.Type,
			&entry.ExtraData,
			&entry.CreatedAt,
			&entry.UserID,
		)

		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
