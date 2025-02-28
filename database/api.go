package database

import (
	"fmt"
	"net/url"
	"strings"

	postgres_db "github.com/papermerge/pmdump/database/postgres"
	sqlite_db "github.com/papermerge/pmdump/database/sqlite"
	"github.com/papermerge/pmdump/models"
	"github.com/papermerge/pmdump/types"
)

func Open(dburl string, appVer types.AppVersion) (*types.DBConn, error) {
	parsedDBURL, err := url.Parse(dburl)

	if err != nil {
		return nil, fmt.Errorf("Error parsing dburl %s: %v", dburl, err)
	}

	if strings.HasPrefix(parsedDBURL.Scheme, "sqlite") {
		db, err := sqlite_db.Open(parsedDBURL.Path)
		if err != nil {
			return nil, err
		}
		dbconn := types.DBConn{
			AppVersion: appVer,
			DBType:     types.SQLite,
			DB:         db,
		}

		return &dbconn, nil
	}

	if parsedDBURL.Scheme == "postgres" {
		db, err := postgres_db.Open(dburl)
		if err != nil {
			return nil, err
		}
		dbconn := types.DBConn{
			AppVersion: appVer,
			DBType:     types.Postgres,
			DB:         db,
		}

		return &dbconn, nil
	}

	return nil, fmt.Errorf("unsupported schema %s in %s", parsedDBURL.Scheme, dburl)
}

func GetUsers(db *types.DBConn) ([]models.User, error) {
	if db.DBType == types.SQLite {
		return sqlite_db.GetUsers(db.DB, db.AppVersion)
	}

	if db.DBType == types.Postgres {
		return postgres_db.GetUsers(db.DB, db.AppVersion)
	}

	return nil, fmt.Errorf("GetUsers: DBType %s not supported", db.DBType)
}

func GetHomeFlatNodes(db *types.DBConn, user_id interface{}) ([]models.FlatNode, error) {
	if db.DBType == types.SQLite {
		return sqlite_db.GetHomeFlatNodes(db.DB, db.AppVersion, user_id)
	}

	if db.DBType == types.Postgres {
		return postgres_db.GetHomeFlatNodes(db.DB, db.AppVersion, user_id)
	}

	return nil, fmt.Errorf("GetHomeFlatNodes: DBType %s not supported", db.DBType)
}

func GetInboxFlatNodes(db *types.DBConn, user_id interface{}) ([]models.FlatNode, error) {
	if db.DBType == types.SQLite {
		return sqlite_db.GetHomeFlatNodes(db.DB, db.AppVersion, user_id)
	}

	if db.DBType == types.Postgres {
		return postgres_db.GetInboxFlatNodes(db.DB, db.AppVersion, user_id)
	}

	return nil, fmt.Errorf("GetInboxFlatNodes: DBType %s not supported", db.DBType)
}

func GetUserNodes(db *types.DBConn, user *models.User) error {
	if db.DBType == types.SQLite {
		return sqlite_db.GetUserNodes(db.DB, db.AppVersion, user)
	}

	if db.DBType == types.Postgres {
		return postgres_db.GetUserNodes(db.DB, db.AppVersion, user)
	}

	return fmt.Errorf("GetUserNodes: DBType %s not supported", db.DBType)
}

func GetDocumentPageRows(db *types.DBConn, user_id interface{}) ([]models.DocumentPageRow, error) {
	if db.DBType == types.SQLite {
		return sqlite_db.GetDocumentPageRows(db.DB, db.AppVersion, user_id)
	}

	if db.DBType == types.Postgres {
		return postgres_db.GetDocumentPageRows(db.DB, db.AppVersion, user_id)
	}

	return nil, fmt.Errorf("GetUserNodes: DBType %s not supported", db.DBType)
}
