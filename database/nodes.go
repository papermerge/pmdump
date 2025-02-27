package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/papermerge/pmdump/constants"
	"github.com/papermerge/pmdump/models"
)

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
				n.NodeUUID,
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
	noHyphenID := strings.ReplaceAll(page.UUID.String(), "-", "")
	noHyphenDocumentVersionID := strings.ReplaceAll(docVer.UUID.String(), "-", "")

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
	noHyphenID := strings.ReplaceAll(docVer.UUID.String(), "-", "")
	noHyphenDocumentID := strings.ReplaceAll(n.NodeUUID.String(), "-", "")

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
	noHyphenID := strings.ReplaceAll(n.NodeUUID.String(), "-", "")
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
	noHyphenID := strings.ReplaceAll(n.NodeUUID.String(), "-", "")
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
