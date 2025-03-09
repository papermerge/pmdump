package sqlite_app_v3_2

import (
	"database/sql"
	"fmt"
	"os"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/utils"
)

const (
	inboxLen = len(".inbox")
	homeLen  = len(".home")
)

/* Strips "inbox" prefix from the string */
func WithoutInboxPrefix(path string) string {
	return path[inboxLen:]
}

/* Strips "home" prefix from the string */
func WithoutHomePrefix(path string) string {
	return path[homeLen:]
}

func GetInboxFlatNodes(db *sql.DB, user_id interface{}) ([]models.FlatNode, error) {
	query := `
    WITH RECURSIVE node_tree AS (
      SELECT
        n.id,
        n.title,
        n.ctype AS model,
        n.title as full_path
      FROM core_basetreenode n
      WHERE parent_id is NULL AND title = '.inbox' AND user_id = ?

      UNION ALL

      SELECT
        n.id,
        n.title,
        n.ctype AS MODEL,
        CAST(CONCAT_WS('/', nt.full_path, n.title) AS VARCHAR(200)) AS full_path
      FROM core_basetreenode n
      INNER JOIN node_tree nt ON n.parent_id = nt.id
      LEFT JOIN core_document doc ON doc.basetreenode_ptr_id = n.id
      WHERE n.user_id = ?
    )
    SELECT
      id,
      title,
      model,
      full_path,
      LENGTH(full_path) AS path_len
    FROM node_tree
    ORDER BY path_len ASC;
  `
	user_uuid := utils.AnyUUID2STR(user_id)

	rows, err := db.Query(query, user_uuid, user_uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []models.FlatNode
	var discard int

	for rows.Next() {
		var node models.FlatNode
		err = rows.Scan(
			&node.ID,
			&node.Title,
			&node.Model,
			&node.FullPath,
			&discard,
		)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}
	return nodes, nil
}

func GetHomeFlatNodes(db *sql.DB, user_id interface{}) ([]models.FlatNode, error) {
	query := `
    WITH RECURSIVE node_tree AS (
      SELECT
        n.id,
        n.title,
        n.ctype AS model,
        n.title as full_path
      FROM core_basetreenode n
      WHERE parent_id is NULL AND title = '.home' AND user_id = ?

      UNION ALL

      SELECT
        n.id,
        n.title,
        n.ctype AS MODEL,
        CAST(CONCAT_WS('/', nt.full_path, n.title) AS VARCHAR(200)) as full_path
      FROM core_basetreenode n
      INNER JOIN node_tree nt ON n.parent_id = nt.id
      LEFT JOIN core_document doc ON doc.basetreenode_ptr_id = n.id
      WHERE n.user_id = ?
    )
    SELECT
      id,
      title,
      model,
      full_path,
      LENGTH(full_path) AS path_len
    FROM node_tree
    ORDER BY path_len ASC;
  `
	user_uuid := utils.AnyUUID2STR(user_id)

	rows, err := db.Query(query, user_uuid, user_uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []models.FlatNode
	var discard int

	for rows.Next() {
		var node models.FlatNode
		err = rows.Scan(
			&node.ID,
			&node.Title,
			&node.Model,
			&node.FullPath,
			&discard,
		)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}
	return nodes, nil
}

func GetUserNodes(db *sql.DB, u *interface{}) error {

	user := (*u).(*models.User)

	user.Inbox = &models.Node{
		Title:    "inbox",
		ID:       user.InboxFolderID,
		NodeType: models.NodeFolderType,
	}
	user.Home = &models.Node{
		Title:    "home",
		ID:       user.HomeFolderID,
		NodeType: models.NodeFolderType,
	}
	homeFlatNodes, err := GetHomeFlatNodes(db, user.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in GetHomeFlatNodes: %v\n", err)
		os.Exit(1)
	}

	for _, node := range homeFlatNodes {
		if node.FullPath == ".home" {
			continue
		}
		node.FullPath = WithoutHomePrefix(node.FullPath)
		user.Home.Insert(node)
	}

	inboxFlatNodes, err := GetInboxFlatNodes(db, user.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in GetInboxFlatNodes: %v\n", err)
		os.Exit(1)
	}

	for _, node := range inboxFlatNodes {
		if node.FullPath == ".inbox" {
			continue
		}
		node.FullPath = WithoutInboxPrefix(node.FullPath)
		user.Inbox.Insert(node)
	}

	return nil
}
