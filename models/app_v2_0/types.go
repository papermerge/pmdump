package models_app_v2_0

import (
	"database/sql"

	"github.com/google/uuid"
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
	LegacyID int `yaml:"legacy_id"`
	Username string
	EMail    string
	Home     *Node
	Inbox    *Node
}

type Users []User

type FlatNode struct {
	ID        int
	Title     string
	Model     string
	FullPath  string
	FileName  *string
	PageCount *int
	Version   *int
}

type Node struct {
	ID        uuid.UUID
	LegacyID  int               `yaml:"legacy_id"`
	Title     string            `yaml:"title"`
	Children  map[string]*Node  `yaml:"children,omitempty"`
	NodeType  NodeType          `yaml:"node_type,omitempty"`
	Versions  []DocumentVersion `yaml:"versions,omitempty"`
	FileName  *string           `yaml:"file_name,omitempty"`
	PageCount *int              `yaml:"page_count,omitempty"`
	Version   *int              `yaml:"version,omitempty"`
}

type DocumentVersion struct {
	ID       uuid.UUID
	LegacyID int `yaml:"legacy_id"`
	Number   int
	FileName string `yaml:"file_name"`
	Pages    []Page
}

type Page struct {
	ID       uuid.UUID
	LegacyID int `yaml:"legacy_id"`
	Text     string
	Number   int
}

type Data struct {
	Users []User
}

type DocumentPageRow struct {
	PageLegacyID     int
	PageID           uuid.UUID
	PageNumber       int
	Text             string
	DocumentID       uuid.UUID
	DocumentLegacyID int
	DocumentVersion  int
}

type DocumentPageRows []DocumentPageRow

type NodeOperation func(n *Node, user_id int, docPages []DocumentPageRow, mediaRoot string)
type NodeQuickOperation func(n *Node)

type TargetNodeOperation func(db *sql.DB, userID uuid.UUID, rootID uuid.UUID, source *Node)
