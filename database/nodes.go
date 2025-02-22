package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/papermerge/pmdump/models"
)

func GetInboxFlatNodes(db *sql.DB, user_id int) ([]models.FlatNode, error) {
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

func GetHomeFlatNodes(db *sql.DB, user_id int) ([]models.FlatNode, error) {
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

func ForEachNode(
	db *sql.DB,
	n *models.Node,
	targetRootID uuid.UUID,
	targetUserID uuid.UUID,
	op models.TargetNodeOperation,
) {

	if n.NodeType == models.DocumentType {
		InsertDocument(db, n, targetRootID, targetUserID)
	} else {
		if err := InsertFolder(db, n, targetRootID, targetUserID); err != nil {
			fmt.Fprintf(os.Stderr, "Node operation error: %v\n", err)
		}
	}

	for _, child := range n.Children {
		ForEachNode(
			db,
			child,
			targetRootID,
			targetUserID,
			op,
		)
	}
}

func InsertDocument(
	db *sql.DB,
	n *models.Node,
	targetRootID uuid.UUID,
	targetUserID uuid.UUID,
) {

}

func InsertFolder(
	db *sql.DB,
	n *models.Node,
	parentID uuid.UUID,
	userID uuid.UUID,
) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Defer a rollback in case of failure
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	id, _ := uuid.NewUUID()
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	result, err := db.Exec(
		"INSERT INTO nodes (id, title, lang, ctype, user_id, parent_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		id, n.Title, "eng", "folder", userID, parentID, currentTime, currentTime,
	)
	if err != nil {
		return fmt.Errorf(
			"insert node %q, parentID %q, userID %q: %v",
			n.Title,
			parentID,
			userID,
			err,
		)
	}

	// Get the last inserted ID
	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last inserted ID %s: %v", n.Title, err)
	}

	_, err = db.Exec(
		"INSERT INTO folders (node_id) VALUES (?)",
		lastInsertedID,
	)
	if err != nil {
		return fmt.Errorf("insert folder %s: %v", n.Title, err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction for %s: %v", n.Title, err)
	}

	return nil
}

func CreateTargetNode(db *sql.DB, userID uuid.UUID, rootID uuid.UUID, source *models.Node) {

}
