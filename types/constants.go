package types

const (
	V2_0 AppVersion = "2.0"
	V2_1 AppVersion = "2.1"
	V3_0 AppVersion = "3.0"
	V3_1 AppVersion = "3.1"
	V3_2 AppVersion = "3.2"
	V3_3 AppVersion = "3.3"
	V3_4 AppVersion = "3.4"
)

var AppVersionsForExport = []AppVersion{
	V2_0,
	V2_1,
	V3_0,
	V3_1,
	V3_2,
	V3_3,
	V3_4,
}

const (
	SQLite   DBType = "sqlite"
	Postgres DBType = "postgres"
)
