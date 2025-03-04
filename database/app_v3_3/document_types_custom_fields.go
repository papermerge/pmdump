package database_app_v3_3

import (
	_ "github.com/mattn/go-sqlite3"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/types"
)

func GetDocumentTypesCustomFields(db *types.DBConn) ([]models.DocumentTypesCustomFields, error) {
	query := `
    SELECT
      id,
      document_type_id,
      custom_field_id
    FROM document_types_custom_fields
  `
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.DocumentTypesCustomFields

	for rows.Next() {
		var entry models.DocumentTypesCustomFields

		err = rows.Scan(
			&entry.ID,
			&entry.DocumentTypeID,
			&entry.CustomFieldID,
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
