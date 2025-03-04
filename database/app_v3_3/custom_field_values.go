package database_app_v3_3

import (
	_ "github.com/mattn/go-sqlite3"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/types"
)

func GetCustomFieldValues(db *types.DBConn) ([]models.CustomFieldValues, error) {
	query := `
    SELECT
      id,
      document_id,
      field_id,
      value_text,
      value_boolean,
      value_date,
      value_int,
      value_float,
      value_monetary,
      value_yearmonth,
      created_at
    FROM custom_field_values
  `
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.CustomFieldValues

	for rows.Next() {
		var entry models.CustomFieldValues

		err = rows.Scan(
			&entry.ID,
			&entry.DocumentID,
			&entry.FieldID,
			&entry.ValueText,
			&entry.ValueBoolean,
			&entry.ValueDate,
			&entry.ValueInt,
			&entry.ValueFloat,
			&entry.ValueMonetary,
			&entry.ValueYearMonth,
			&entry.CreatedAt,
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
