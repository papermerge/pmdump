package models_app_v3_3

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/papermerge/pmdump/types"
)

type NodeType string

const (
	NodeFolderType   NodeType = "folder"
	NodeDocumentType NodeType = "document"
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
	ID            uuid.UUID
	HomeFolderID  uuid.UUID `yaml:"-"` // this field will be skipped
	InboxFolderID uuid.UUID `yaml:"-"` // this field will be skipped
	Username      string
	EMail         string
	Home          *Node
	Inbox         *Node
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
	Number   int
	FileName string `yaml:"file_name"`
	Size     int
	Lang     string
	Text     *string `yaml:"text,omitempty"`
	Pages    []Page
}

type Page struct {
	ID     uuid.UUID
	Text   *string `yaml:"text,omitempty"`
	Number int
}

type Data struct {
	Users                     []User
	Groups                    []Group
	Permissions               []Permission
	GroupsPermissions         []GroupsPermissions `yaml:"groups_permissions"`
	DocumentTypes             []DocumentType      `yaml:"document_types"`
	Tags                      []Tag
	NodesTags                 []NodesTags                 `yaml:"nodes_tags"`
	UsersGroups               []UsersGroups               `yaml:"users_groups"`
	UsersPermissions          []UsersPermissions          `yaml:"users_permissions"`
	CustomFields              []CustomField               `yaml:"custom_fields"`
	DocumentTypesCustomFields []DocumentTypesCustomFields `yaml:"document_types_custom_fields"`
	CustomFieldValues         []CustomFieldValues         `yaml:"custom_field_values"`
}

type DocumentVersionPageRow struct {
	PageID                uuid.UUID
	PageNumber            int
	PageText              *string
	FileName              string
	Size                  int
	Lang                  string
	DocumentID            uuid.UUID
	DocumentVersionID     uuid.UUID
	DocumentVersionText   *string
	DocumentVersionNumber int
}

type NodeOperation func(db *types.DBConn, n any) error
type NodeQuickOperation func(n *Node)

type TargetNodeOperation func(db *sql.DB, userID uuid.UUID, rootID uuid.UUID, source *Node)

type Group struct {
	ID   uuid.UUID
	Name string
}

type Permission struct {
	ID       uuid.UUID
	Name     string
	Codename string
}

type GroupsPermissions struct {
	GroupID      uuid.UUID `yaml:"group_id"`
	PermissionID uuid.UUID `yaml:"permission_id"`
}

type DocumentType struct {
	ID           uuid.UUID
	Name         string
	PathTemplate string    `yaml:"path_template"`
	UserID       uuid.UUID `yaml:"user_id"`
	CreatedAt    time.Time `yaml:"created_at"`
}

type Tag struct {
	ID          uuid.UUID
	Name        string
	FGColor     string `yaml:fg_color"`
	BGColor     string `yaml:bg_color"`
	Pinned      bool
	Description string
	UserID      string `yaml:"user_id"`
}

type NodesTags struct {
	ID     int
	NodeID uuid.UUID `yaml:"node_id"`
	TagID  uuid.UUID `yaml:"tag_id"`
}

type UsersGroups struct {
	UserID  uuid.UUID `yaml:"user_id"`
	GroupID uuid.UUID `yaml:"group_id"`
}

type UsersPermissions struct {
	UserID       uuid.UUID `yaml:"user_id"`
	PermissionID uuid.UUID `yaml:"permission_id"`
}

type CustomField struct {
	ID        uuid.UUID
	Name      string
	Type      string
	ExtraData *string   `yaml:"extra_data,omitempty"`
	CreatedAt time.Time `yaml:"created_at"`
	UserID    uuid.UUID `yaml:"user_id"`
}

type DocumentTypesCustomFields struct {
	ID             int32
	DocumentTypeID uuid.UUID `yaml:"document_type_id"`
	CustomFieldID  uuid.UUID `yaml:"custom_field_id"`
}

type CustomFieldValues struct {
	ID             uuid.UUID
	DocumentID     uuid.UUID  `yaml:"document_id"`
	FieldID        uuid.UUID  `yaml:"field_id"`
	ValueText      *string    `yaml:"value_text,omitempty"`
	ValueBoolean   *bool      `yaml:"value_boolean,omitempty"`
	ValueDate      *time.Time `yaml:"value_date,omitempty"`
	ValueInt       *int32     `yaml:"value_int,omitempty"`
	ValueFloat     *float32   `yaml:"value_float,omitempty"`
	ValueMonetary  *float32   `yaml:"value_monetary,omitempty"`
	ValueYearMonth *float32   `yaml:"value_yearmonth,omitempty"`
	CreatedAt      time.Time  `yaml:"created_at"`
}
