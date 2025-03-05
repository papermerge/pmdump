package postgres_app_v3_4

import (
	"database/sql"
	"fmt"

	models "github.com/papermerge/pmdump/models/app_v3_3"
)

func InsertDocumentTypesCustomFields(db *sql.DB, dt any) error {
	dtcfs := dt.([]models.DocumentTypesCustomFields)
	query := `
    INSERT INTO
      document_types_custom_fields (id, document_type_id, custom_field_id)
      VALUES(?, ?, ?)
    `

	for _, dtcf := range dtcfs {
		_, err := db.Exec(query,
			dtcf.ID,
			dtcf.DocumentTypeID,
			dtcf.CustomFieldID,
		)

		if err != nil {
			return fmt.Errorf(
				"insert document_types_custom_fields %v failed: %v",
				dtcf,
				err,
			)
		}
	}

	return nil
}
