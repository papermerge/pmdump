package types

import (
	"database/sql"

	"github.com/google/uuid"
)

type AppVersion string
type DBType string

type DBConn struct {
	DB         *sql.DB
	AppVersion AppVersion
	DBType     DBType
}

type FilePath struct {
	Source string
	Dest   string
}

type UserIDChange struct {
	SourceUserID uuid.UUID
	TargetUserID uuid.UUID
}
