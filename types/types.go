package types

import "database/sql"

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
