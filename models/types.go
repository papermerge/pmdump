package models

const (
	FolderModelName   = "folder"
	DocumentModelName = "document"
)

type User struct {
	ID       int
	Username string
	EMail    string
}

type Node struct {
	ID        int
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
	Title    string
	UserID   int  `yaml:"user_id"`
	ParentID *int `yaml:"parent_id"`
}

type Document struct {
	ID        int
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
