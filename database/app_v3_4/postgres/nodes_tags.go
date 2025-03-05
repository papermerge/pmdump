package postgres_app_v3_4

import (
	"database/sql"
	"fmt"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/utils"
)

func InsertNodesTags(db *sql.DB, nt any) error {
	nodes_tags := nt.([]models.NodesTags)
	query := `INSERT INTO nodes_tags (id, node_id, tag_id) VALUES(?, ?, ?)`

	for _, node_tag := range nodes_tags {
		id := node_tag.ID
		node_id := utils.UUID2STR(node_tag.NodeID)
		tag_id := utils.UUID2STR(node_tag.TagID)
		_, err := db.Exec(
			query,
			id,
			node_id,
			tag_id,
		)

		if err != nil {
			return fmt.Errorf(
				"insert node_tag %v failed: %v",
				node_tag,
				err,
			)
		}
	}

	return nil
}
