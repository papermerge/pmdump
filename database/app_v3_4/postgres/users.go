package postgres_app_v3_4

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/papermerge/pmdump/constants"
	models "github.com/papermerge/pmdump/models/app_v3_3"
)

func GetTargetUsers(db *sql.DB) (models.TargetUserList, error) {
	rows, err := db.Query("SELECT id, username, email, home_folder_id, inbox_folder_id FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users models.TargetUserList

	for rows.Next() {
		var user models.TargetUser
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.EMail,
			&user.HomeID,
			&user.InboxID,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func InsertUsersData(
	db *sql.DB,
	su interface{},
	tu interface{},
) {

	sourceUsers := su.(models.Users)
	targetUsers := tu.(models.TargetUserList)

	for i := 0; i < len(sourceUsers); i++ {
		targetUser := targetUsers.Get(sourceUsers[i].Username)
		if targetUser != nil {
			ImportUserData(db, sourceUsers[i], targetUser)
		} else {
			targetUser, err := CreateTargetUser(db, sourceUsers[i])
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"Error creating target user for %s: %v\n",
					sourceUsers[i].Username,
					err,
				)
				continue
			}
			ImportUserData(db, sourceUsers[i], targetUser)
		}
	}
}

func ImportUserData(
	db *sql.DB,
	sourceUser models.User,
	targetUser *models.TargetUser,
) {

	models.ForEachNode(sourceUser.Home, models.UpdateNodeUUID)
	models.ForEachNode(sourceUser.Inbox, models.UpdateNodeUUID)

	ForEachSourceNode(
		db,
		sourceUser.Home, // start here
		targetUser.HomeID,
		targetUser.ID,
		CreateTargetNode,
	)
	ForEachSourceNode(
		db,
		sourceUser.Inbox, // start here
		targetUser.InboxID,
		targetUser.ID,
		CreateTargetNode,
	)
}

func CreateTargetUser(
	db *sql.DB,
	source models.User,
) (*models.TargetUser, error) {
	return nil, nil
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
