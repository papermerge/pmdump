package sqlite_app_v3_3

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/papermerge/pmdump/constants"
	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/utils"
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
	// works only for sqlite3 (because of "||"... For PostgreSQL use "concat")
	query := `
    WITH RECURSIVE node_tree AS (
      SELECT
        n.id,
        n.title,
        n.ctype AS model,
        n.title as full_path
      FROM nodes n
      WHERE parent_id is NULL AND title = 'home' AND user_id = ?

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
		if node.FullPath == "inbox" {
			continue
		}
		node.FullPath = utils.WithoutInboxPrefix(node.FullPath)
		user.Inbox.Insert(node)
	}

	return nil
}

func ForEachSourceNode(
	db *sql.DB,
	n *models.Node,
	targetParentID uuid.UUID,
	targetUserID uuid.UUID,
	op models.TargetNodeOperation,
) {

	if n.NodeType == models.DocumentType {
		if err := InsertDocument(db, n, targetParentID, targetUserID); err != nil {
			fmt.Fprintf(os.Stderr, "Document insert error: %v\n", err)
		}
	} else {
		if n.Title != constants.INBOX && n.Title != constants.HOME {
			if err := InsertFolder(db, n, targetParentID, targetUserID); err != nil {
				fmt.Fprintf(os.Stderr, "Folder insert error: %v\n", err)
			}
		}
	}

	for _, child := range n.Children {
		if n.Title == constants.HOME || n.Title == constants.INBOX {
			ForEachSourceNode(
				db,
				child,
				// this is either Home ID or Inbox ID of target user
				targetParentID,
				targetUserID,
				op,
			)
		} else {
			ForEachSourceNode(
				db,
				child,
				n.ID,
				targetUserID,
				op,
			)
		}
	}
}

func InsertPage(
	db *sql.DB,
	docVer models.DocumentVersion,
	page models.Page,
) error {
	noHyphenID := strings.ReplaceAll(page.ID.String(), "-", "")
	noHyphenDocumentVersionID := strings.ReplaceAll(docVer.ID.String(), "-", "")

	_, err := db.Exec(
		"INSERT INTO pages (id, document_version_id, number, page_count, lang) VALUES (?, ?, ?, ?, ?)",
		noHyphenID,
		noHyphenDocumentVersionID,
		page.Number,
		len(docVer.Pages),
		constants.ENG,
	)

	if err != nil {
		return fmt.Errorf(
			"insert page ID=%q, number %q failed: %v",
			noHyphenID,
			page.Number,
			err,
		)
	}

	return nil
}

func InsertDocumentVersion(
	db *sql.DB,
	n *models.Node,
	docVer models.DocumentVersion,
) error {
	noHyphenID := strings.ReplaceAll(docVer.ID.String(), "-", "")
	noHyphenDocumentID := strings.ReplaceAll(n.ID.String(), "-", "")

	_, err := db.Exec(
		"INSERT INTO document_versions (id, document_id, number, file_name, lang, size, page_count) VALUES (?, ?, ?, ?, ?, ?, ?)",
		noHyphenID,
		noHyphenDocumentID,
		docVer.Number,
		docVer.FileName,
		constants.ENG,
		0,
		len(docVer.Pages),
	)
	if err != nil {
		return fmt.Errorf(
			"insert document version %q, number %q, file_name %q failed: %v",
			noHyphenID,
			docVer.Number,
			docVer.FileName,
			err,
		)
	}

	for _, page := range docVer.Pages {
		err = InsertPage(
			db, docVer, page,
		)
		if err != nil {
			return fmt.Errorf(
				"insert page for document version %q, number %q, file_name %q failed: %v",
				noHyphenID,
				docVer.Number,
				docVer.FileName,
				err,
			)
		}
	}

	return nil
}

func InsertDocument(
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

	noHyphenParentID := strings.ReplaceAll(parentID.String(), "-", "")
	noHyphenID := strings.ReplaceAll(n.ID.String(), "-", "")
	noHyphenUserID := strings.ReplaceAll(userID.String(), "-", "")

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	_, err = db.Exec(
		"INSERT INTO nodes (id, title, lang, ctype, user_id, parent_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		noHyphenID,
		n.Title,
		constants.ENG,
		constants.DOCUMENT,
		noHyphenUserID,
		noHyphenParentID,
		currentTime,
		currentTime,
	)
	if err != nil {
		return fmt.Errorf(
			"insert node %q, parentID %q, userID %q: %v",
			n.Title,
			noHyphenParentID,
			noHyphenUserID,
			err,
		)
	}

	_, err = db.Exec(
		"INSERT INTO documents (node_id, ocr, ocr_status) VALUES (?, ?, ?)",
		noHyphenID,
		false,
		constants.UNKNOWN,
	)
	if err != nil {
		return fmt.Errorf("insert document %s: %v", n.Title, err)
	}

	for _, docVer := range n.Versions {
		err = InsertDocumentVersion(
			db, n, docVer,
		)
		if err != nil {
			return fmt.Errorf(
				"insert document %q, documentID %q, parentID %q, userID %q failed: %v",
				n.Title,
				noHyphenID,
				noHyphenParentID,
				noHyphenUserID,
				err,
			)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction for %s: %v", n.Title, err)
	}
	return nil
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

	noHyphenParentID := strings.ReplaceAll(parentID.String(), "-", "")
	noHyphenID := strings.ReplaceAll(n.ID.String(), "-", "")
	noHyphenUserID := strings.ReplaceAll(userID.String(), "-", "")

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	_, err = db.Exec(
		"INSERT INTO nodes (id, title, lang, ctype, user_id, parent_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		noHyphenID,
		n.Title,
		constants.ENG,
		constants.FOLDER,
		noHyphenUserID,
		noHyphenParentID,
		currentTime,
		currentTime,
	)
	if err != nil {
		return fmt.Errorf(
			"insert node %q, parentID %q, userID %q: %v",
			n.Title,
			noHyphenParentID,
			noHyphenUserID,
			err,
		)
	}
	_, err = db.Exec(
		"INSERT INTO folders (node_id) VALUES (?)",
		noHyphenID,
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
