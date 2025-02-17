package models

import (
	"github.com/google/uuid"
)

const (
	FolderModelName   = "folder"
	DocumentModelName = "document"
)

type ID2UUID map[int]uuid.UUID

type IDDict struct {
	NodeIDs ID2UUID
	UserIDs ID2UUID
}

type User struct {
	ID       int
	UUID     uuid.UUID
	Username string
	EMail    string
}

type Node struct {
	ID         int
	UUID       uuid.UUID
	Title      string
	Model      string
	UserID     int       `yaml:"user_id"`
	FileName   *string   `yaml:"file_name"`
	PageCount  *int      `yaml:"page_count"`
	ParentID   *int      `yaml:"parent_id"`
	ParentUUID uuid.UUID `yaml:"parent_uuid"`
	Version    *int
}

type Folder struct {
	ID         int
	UUID       uuid.UUID
	Title      string
	UserID     int        `yaml:"user_id"`
	UserUUID   uuid.UUID  `yaml:"user_uuid"`
	ParentID   *int       `yaml:"parent_id"`
	ParentUUID *uuid.UUID `yaml:"parent_uuid"`
}

type DocumentVersion struct {
	Number    int
	UUID      uuid.UUID
	FileName  *string `yaml:"file_name"`
	PageCount *int    `yaml:"page_count"`
	Pages     []Page
}

type Document struct {
	ID         int
	UUID       uuid.UUID
	Title      string
	UserID     int        `yaml:"user_id"`
	ParentID   *int       `yaml:"parent_id"`
	UserUUID   uuid.UUID  `yaml:"user_uuid"`
	ParentUUID *uuid.UUID `yaml:"parent_uuid"`
	Versions   []DocumentVersion
}

type Page struct {
	ID     int
	UUID   uuid.UUID
	Number int
	Text   string
}

type Data struct {
	Users     []User
	Documents []Document
	Folders   []Folder
	Tags      []Tag
}

type FilePath struct {
	Source string
	Dest   string
}

type DocumentPageRow struct {
	PageID          int
	PageUUID        uuid.UUID
	PageNumber      int
	Text            string
	DocumentID      int
	DocumentVersion int
}

type Tag struct {
	ID          int
	UUID        uuid.UUID
	Name        string
	Description string
	BGColor     string
	FGColor     string
	Pinned      bool
}
