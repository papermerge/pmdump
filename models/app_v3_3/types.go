package models_app_v3_3

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/papermerge/pmdump/types"
)

type NodeType string

const (
	FolderType   NodeType = "folder"
	DocumentType NodeType = "document"
)

/*
User as read from target database
Target = Destination = Papermerge 3.4
*/
type TargetUser struct {
	ID       uuid.UUID
	Username string
	EMail    string
	HomeID   uuid.UUID
	InboxID  uuid.UUID
}

type TargetUserList []TargetUser

type BaseUser struct {
	Username string
	EMail    string
	Home     *Node
	Inbox    *Node
}

type User struct {
	ID       uuid.UUID
	Username string
	EMail    string
	Home     *Node
	Inbox    *Node
}

type Users []User

type FlatNode struct {
	ID        uuid.UUID
	Title     string
	Model     string
	FullPath  string
	FileName  *string
	PageCount *int
	Version   *int
}

type Node struct {
	ID        uuid.UUID
	Title     string           `yaml:"title"`
	Children  map[string]*Node `yaml:"children,omitempty"`
	NodeType  NodeType
	Versions  []DocumentVersion `yaml:"versions,omitempty"`
	FileName  *string           `yaml:"file_name,omitempty"`
	PageCount *int              `yaml:"page_count,omitempty"`
	Version   *int              `yaml:"version,omitempty"`
}

type DocumentVersion struct {
	ID       uuid.UUID
	Number   int
	FileName string `yaml:"file_name"`
	Pages    []Page
}

type Page struct {
	ID     uuid.UUID
	Text   *string
	Number int
}

type Data struct {
	Users []User
}

type DocumentVersionPageRow struct {
	PageID                uuid.UUID
	PageNumber            int
	PageText              *string
	DocumentID            uuid.UUID
	DocumentVersionID     uuid.UUID
	DocumentVersionText   *string
	DocumentVersionNumber int
}

type NodeOperation func(db *types.DBConn, n any) error
type NodeQuickOperation func(n *Node)

type TargetNodeOperation func(db *sql.DB, userID uuid.UUID, rootID uuid.UUID, source *Node)
