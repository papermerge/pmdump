package models

import (
	"github.com/google/uuid"
)

const (
	FolderModelName   = "folder"
	DocumentModelName = "document"
)

type User struct {
	ID       int
	UUID     uuid.UUID
	Username string
	EMail    string
}

type Node struct {
	ID        int
	UUID      uuid.UUID
	Title     string
	Model     string
	UserID    int  `yaml:"user_id"`
	ParentID  *int `yaml:"parent_id"`
	Version   *int
	FileName  *string `yaml:"file_name"`
	PageCount *int    `yaml:"page_count"`
}

type Folder struct {
	ID       int
	UUID     uuid.UUID
	Title    string
	UserID   int  `yaml:"user_id"`
	ParentID *int `yaml:"parent_id"`
}

type Document struct {
	ID        int
	UUID      uuid.UUID
	Title     string
	UserID    int  `yaml:"user_id"`
	ParentID  *int `yaml:"parent_id"`
	Version   *int
	FileName  *string `yaml:"file_name"`
	PageCount *int    `yaml:"page_count"`
}

type Data struct {
	Users     []User
	Documents []Document
	Folders   []Folder
}

type FilePath struct {
	Source string
	Dest   string
}
