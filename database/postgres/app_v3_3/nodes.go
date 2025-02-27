package postgres_app_v3_3

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
        ct.model AS model,
        n.title as full_path,
        doc.version,
        doc.file_name,
        doc.page_count
      FROM core_basetreenode n
      INNER JOIN django_content_type ct ON ct.id = n.polymorphic_ctype_id
      LEFT JOIN core_document doc ON doc.basetreenode_ptr_id = n.id
      WHERE parent_id is NULL and title = '.inbox' AND user_id = ?

      UNION ALL

      SELECT
        n.id,
        n.title,
        ct.model AS model,
        nt.full_path || '/' || n.title AS full_path,
        doc.version,
        doc.file_name,
        doc.page_count
      FROM core_basetreenode n
      INNER JOIN node_tree nt ON n.parent_id = nt.id
      INNER JOIN django_content_type ct ON ct.id = n.polymorphic_ctype_id
      LEFT JOIN core_document doc ON doc.basetreenode_ptr_id = n.id
      WHERE n.user_id = ?
    )
    SELECT
      id,
      title,
      model,
      full_path,
      LENGTH(full_path) AS path_len,
      version,
      file_name,
      page_count
    FROM node_tree
    ORDER BY path_len ASC;
  `
	rows, err := db.Query(query, user_id, user_id)
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
			&node.Version,
			&node.FileName,
			&node.PageCount,
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
        ct.model AS model,
        n.title as full_path,
        doc.version,
        doc.file_name,
        doc.page_count
      FROM core_basetreenode n
      INNER JOIN django_content_type ct ON ct.id = n.polymorphic_ctype_id
      LEFT JOIN core_document doc ON doc.basetreenode_ptr_id = n.id
      WHERE parent_id is NULL AND title != '.inbox' AND user_id = ?

      UNION ALL

      SELECT
        n.id,
        n.title,
        ct.model AS MODEL,
        nt.full_path || '/' || n.title AS full_path,
        doc.version,
        doc.file_name,
        doc.page_count
      FROM core_basetreenode n
      INNER JOIN node_tree nt ON n.parent_id = nt.id
      INNER JOIN django_content_type ct ON ct.id = n.polymorphic_ctype_id
      LEFT JOIN core_document doc ON doc.basetreenode_ptr_id = n.id
      WHERE n.user_id = ?
    )
    SELECT
      id,
      title,
      model,
      full_path,
      LENGTH(full_path) AS path_len,
      version,
      file_name,
      page_count
    FROM node_tree
    ORDER BY path_len ASC;
  `
	rows, err := db.Query(query, user_id, user_id)
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
			&node.Version,
			&node.FileName,
			&node.PageCount,
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
