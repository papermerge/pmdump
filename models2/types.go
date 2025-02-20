package models2

import "github.com/google/uuid"

type NodeType string

const (
	FolderType   NodeType = "folder"
	DocumentType NodeType = "document"
)

type User struct {
	ID       int
	Username string
	EMail    string
	Home     *Node
	Inbox    *Node
}

type FlatNode struct {
	ID       int
	Title    string
	Model    string
	FullPath string
}

type Node struct {
	ID       int
	Title    string           `yaml:"title"`
	Children map[string]*Node `yaml:"children,omitempty"`
	NodeType NodeType
	Versions []DocumentVersion
}

type DocumentVersion struct {
	UUID     uuid.UUID
	Number   int
	FileName string `yaml:"file_name"`
	Pages    []Page
}

type Page struct {
	UUID uuid.UUID
	Text string
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
