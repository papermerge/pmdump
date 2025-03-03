package database_app_v3_3

import (
	_ "github.com/mattn/go-sqlite3"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/types"
)

func GetDocumentTypes(db *types.DBConn) ([]models.DocumentType, error) {
	rows, err := db.DB.Query("SELECT id, name, path_template, user_id, created_at FROM document_types")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.DocumentType

	for rows.Next() {
		var entry models.DocumentType

		err = rows.Scan(
			&entry.ID,
			&entry.Name,
			&entry.PathTemplate,
			&entry.UserID,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
