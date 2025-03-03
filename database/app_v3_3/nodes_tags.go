package database_app_v3_3

import (
	_ "github.com/mattn/go-sqlite3"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/types"
)

func GetNodesTags(db *types.DBConn) ([]models.NodesTags, error) {
	rows, err := db.DB.Query("SELECT id, node_id, tag_id FROM nodes_tags")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.NodesTags

	for rows.Next() {
		var entry models.NodesTags

		err = rows.Scan(
			&entry.ID,
			&entry.NodeID,
			&entry.TagID,
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
