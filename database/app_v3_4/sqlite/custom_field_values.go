package sqlite_app_v3_4

import (
	"database/sql"
	"fmt"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/utils"
)

func InsertCustomFieldValues(db *sql.DB, cfv any) error {
	custom_field_values := cfv.([]models.CustomFieldValues)
	query := `
    INSERT INTO custom_field_values (
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
    )
    VALUES(
      ?, ?, ?,
      ?, ?, ?,
      ?, ?, ?,
      ?, ?
    )`

	for _, cfv := range custom_field_values {
		id := utils.UUID2STR(cfv.ID)
		document_id := utils.UUID2STR(cfv.DocumentID)
		field_id := utils.UUID2STR(cfv.FieldID)

		_, err := db.Exec(query,
			id,
			document_id,
			field_id,
			cfv.ValueText,
			cfv.ValueBoolean,
			cfv.ValueDate,
			cfv.ValueInt,
			cfv.ValueFloat,
			cfv.ValueMonetary,
			cfv.ValueYearMonth,
			cfv.CreatedAt,
		)

		if err != nil {
			return fmt.Errorf(
				"insert custom_field_value %v failed: %v",
				cfv,
				err,
			)
		}
	}

	return nil
}
