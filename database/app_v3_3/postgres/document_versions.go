package postgres_app_v3_3

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/utils"
)

func GetDocumentVersionsForNode(
	db *sql.DB,
	node_id uuid.UUID,
) ([]models.DocumentVersionPageRow, error) {
	query := `
    SELECT
      d.node_id AS DocumentID,
      dv.id AS DocumentVersionID,
      dv.number AS DocumentVersionNumber,
      dv.text AS DocumentText,
      p.id AS PageID,
      p.number AS PageNumber,
      p.text AS PageText
    FROM document_versions dv
    JOIN pages p ON p.document_version_id = dv.id
    JOIN documents d ON d.node_id = dv.document_id
    WHERE d.node_id = $1
  `
	node_uuid := utils.UUID2STR(node_id)

	rows, err := db.Query(query, node_uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.DocumentVersionPageRow

	for rows.Next() {
		var entry models.DocumentVersionPageRow
		err = rows.Scan(
			&entry.DocumentID,
			&entry.DocumentVersionID,
			&entry.DocumentVersionNumber,
			&entry.DocumentVersionText,
			&entry.PageID,
			&entry.PageNumber,
			&entry.PageText,
		)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}
	return entries, nil
}

func InsertDocVersionsAndPages(
	db *sql.DB,
	n any,
) {
	node := n.(*models.Node)

	docVerPages, err := GetDocumentVersionsForNode(db, node.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	docVersions := make(map[string]models.DocumentVersion)

	for _, docVerEntry := range docVerPages {
		key := docVerEntry.DocumentVersionID.String()

		if docVer, exists := docVersions[key]; exists {
			page := models.Page{
				ID:     docVerEntry.PageID,
				Number: docVerEntry.PageNumber,
				Text:   docVerEntry.PageText,
			}
			docVer.Pages = append(docVer.Pages, page)
		} else {
			page := models.Page{
				ID:     docVerEntry.PageID,
				Number: docVerEntry.PageNumber,
				Text:   docVerEntry.PageText,
			}
			docVer.Pages = append(docVer.Pages, page)
			docVer := models.DocumentVersion{
				ID:       docVerEntry.DocumentVersionID,
				Number:   docVerEntry.DocumentVersionNumber,
				FileName: docVerEntry.FileName,
				Lang:     docVerEntry.Lang,
				Size:     docVerEntry.Size,
				Text:     docVerEntry.DocumentVersionText,
				Pages:    []models.Page{page},
			}
			docVersions[key] = docVer
		}
	}

	for key, _ := range docVersions {
		node.Versions = append(node.Versions, docVersions[key])
	}
}
