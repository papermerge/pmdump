package sqlite_app_v2_0

import (
	"database/sql"

	"github.com/google/uuid"

	_ "github.com/mattn/go-sqlite3"
	models "github.com/papermerge/pmdump/models/app_v2_0"
)

func GetDocumentPageRows(db *sql.DB, user_id interface{}) ([]models.DocumentPageRow, error) {
	query := `
    SELECT p.id,
      p.number,
      p.text,
      p.document_id,
      doc.version
    FROM core_page p
    JOIN core_document doc
      ON p.document_id = doc.basetreenode_ptr_id
    JOIN core_basetreenode node ON node.id = doc.basetreenode_ptr_id
    WHERE node.user_id = ?;
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
			&entry.PageLegacyID,
			&entry.PageNumber,
			&entry.Text,
			&entry.DocumentID,
			&entry.DocumentVersion,
		)
		if err != nil {
			return nil, err
		}
		entry.PageID = uuid.New()
		entries = append(entries, entry)
	}
	return entries, nil
}
