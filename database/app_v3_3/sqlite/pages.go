package sqlite_app_v3_3

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	models "github.com/papermerge/pmdump/models/app_v3_3"
)

func GetDocumentPageRows(db *sql.DB, user_id interface{}) ([]models.DocumentPageRow, error) {
	query := `
    SELECT p.id,
      p.number,
      p.text,
      p.document_version_id,
      doc.node_id AS document_id
    FROM pages p
    JOIN document_versions docver
      ON docver.id = p.document_version_id
    JOIN documents doc
      ON docver.document_id = doc.node_id
    JOIN nodes node
      ON node.id = doc.node_id
    WHERE node.user_id = 'b27994ff777d47eab95b3013a3f3a954';
  `
	rows, err := db.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.DocumentPageRow

	for rows.Next() {
		var entry models.DocumentPageRow
		err = rows.Scan(
			&entry.PageID,
			&entry.PageNumber,
			&entry.Text,
			&entry.DocumentVersionID,
			&entry.DocumentID,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
