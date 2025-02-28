package models

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
	ID       int
	UUID     uuid.UUID
	Username string
	EMail    string
	Home     *Node
	Inbox    *Node
}

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
	ID        int
	NodeUUID  uuid.UUID
	Title     string           `yaml:"title"`
	Children  map[string]*Node `yaml:"children,omitempty"`
	NodeType  NodeType
	Versions  []DocumentVersion `yaml:"versions,omitempty"`
	FileName  *string           `yaml:"file_name,omitempty"`
	PageCount *int              `yaml:"page_count,omitempty"`
	Version   *int              `yaml:"version,omitempty"`
}

type DocumentVersion struct {
	UUID     uuid.UUID
	Number   int
	FileName string `yaml:"file_name"`
	Pages    []Page
}

type Page struct {
	UUID   uuid.UUID
	Text   string
	Number int
}

type Data struct {
	Users []User
}

type DocumentPageRow struct {
	PageID          int
	PageUUID        uuid.UUID
	PageNumber      int
	Text            string
	DocumentID      int
	DocumentVersion int
}

type NodeOperation func(n *Node, user_id int, docPages []DocumentPageRow, mediaRoot string)
type NodeQuickOperation func(n *Node)

type FilePath struct {
	Source string
	Dest   string
}

type TargetNodeOperation func(db *sql.DB, userID uuid.UUID, rootID uuid.UUID, source *Node)
