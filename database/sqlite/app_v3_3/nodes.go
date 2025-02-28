package sqlite_app_v3_3

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/papermerge/pmdump/models"
)

func GetInboxFlatNodes(db *sql.DB, user_id interface{}) ([]models.FlatNode, error) {
	// works only for sqlite3 (because of "||", "?" and "LENGTH". For PostgreSQL use "concat", "$x")
	query := `
    WITH RECURSIVE node_tree AS (
      SELECT
        n.id,
        n.title,
        n.ctype AS model,
        n.title as full_path
      FROM nodes n
      LEFT JOIN documents doc ON doc.node_id = n.id
      WHERE parent_id is NULL and title = 'inbox' AND user_id = ?

      UNION ALL

      SELECT
        n.id,
        n.title,
        n.ctype AS model,
        nt.full_path || '/' || n.title AS full_path
      FROM nodes n
      INNER JOIN node_tree nt ON n.parent_id = nt.id
      LEFT JOIN documents doc ON doc.node_id = n.id
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
	rows, err := db.Query(query, user_id, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []models.FlatNode

	for rows.Next() {
		var node models.FlatNode
		err = rows.Scan(
			&node.ID,
			&node.Title,
			&node.Model,
			&node.FullPath,
		)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}
	return nodes, nil
}

func GetHomeFlatNodes(db *sql.DB, user_id interface{}) ([]models.FlatNode, error) {
	// works only for sqlite3 (because of "||"... For PostgreSQL use "concat")
	query := `
    WITH RECURSIVE node_tree AS (
      SELECT
        n.id,
        n.title,
        n.ctype AS model,
        n.title as full_path
      FROM nodes n
      LEFT JOIN documents doc ON doc.node_id = n.id
      WHERE parent_id is NULL AND title != 'inbox' AND user_id = ?

      UNION ALL

      SELECT
        n.id,
        n.title,
        n.ctype AS MODEL,
        nt.full_path || '/' || n.title AS full_path
      FROM nodes n
      INNER JOIN node_tree nt ON n.parent_id = nt.id
      LEFT JOIN documents doc ON doc.node_id = n.id
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
	rows, err := db.Query(query, user_id, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []models.FlatNode

	for rows.Next() {
		var node models.FlatNode
		err = rows.Scan(
			&node.ID,
			&node.Title,
			&node.Model,
			&node.FullPath,
		)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}
	return nodes, nil
}

func GetUserNodes(db *sql.DB, user *models.User) error {

	user.Inbox = &models.Node{
		Title:    "inbox",
		NodeType: models.FolderType,
	}
	user.Home = &models.Node{
		Title:    "home",
		NodeType: models.FolderType,
	}

	homeFlatNodes, err := GetHomeFlatNodes(db, user.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in GetHomeFlatNodes: %v\n", err)
		os.Exit(1)
	}

	for _, node := range homeFlatNodes {
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
		node.FullPath = node.FullPath[7:]
		user.Inbox.Insert(node)
	}

	return nil
}
