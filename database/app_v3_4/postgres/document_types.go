package postgres_app_v3_4

import (
	"database/sql"
	"fmt"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/utils"
)

func InsertDocumentTypes(db *sql.DB, dt any) error {
	document_types := dt.([]models.DocumentType)
	query := `INSERT INTO document_types (id, name, path_template, user_id, created_at) VALUES($1, $2, $3, $4, $5)`

	for _, document_type := range document_types {
		id := utils.UUID2STR(document_type.ID)
		user_id := utils.UUID2STR(document_type.UserID)
		_, err := db.Exec(query,
			id,
			document_type.Name,
			document_type.PathTemplate,
			user_id,
			document_type.CreatedAt,
		)

		if err != nil {
			return fmt.Errorf(
				"insert document_type %v failed: %v",
				document_type,
				err,
			)
		}
	}

	return nil
}
