package sqlite_db

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	sqlite_app_v2_0 "github.com/papermerge/pmdump/database/sqlite/app_v2_0"
	sqlite_app_v3_3 "github.com/papermerge/pmdump/database/sqlite/app_v3_3"
	"github.com/papermerge/pmdump/models"
	"github.com/papermerge/pmdump/types"
)

func Open(dburl string) (*sql.DB, error) {
	parsedDBURL, err := url.Parse(dburl)

	if err != nil {
		return nil, fmt.Errorf("Error parsing dburl %s: %v", dburl, err)

	}

	if !strings.HasPrefix(parsedDBURL.Scheme, "sqlite") {
		panic("sqlite: schema did not match")
	}

	return sql.Open("sqlite3", dburl)
}

func GetUsers(db *sql.DB, appVer types.AppVersion) ([]models.User, error) {

	if appVer == types.V2_0 {
		return sqlite_app_v2_0.GetUsers(db)
	}
	if appVer == types.V3_3 {
		return sqlite_app_v3_3.GetUsers(db)
	}

	e := fmt.Errorf("GetUsers not implemented for app version %s\n", appVer)
	return nil, e
}

func GetHomeFlatNodes(db *sql.DB, appVer types.AppVersion, user_id interface{}) ([]models.FlatNode, error) {

	if appVer == types.V2_0 {
		return sqlite_app_v2_0.GetHomeFlatNodes(db, user_id)
	}
	if appVer == types.V3_3 {
		return sqlite_app_v3_3.GetHomeFlatNodes(db, user_id)
	}

	e := fmt.Errorf("GetHomeFlatNodes not implemented for app version %s\n", appVer)
	return nil, e
}

func GetInboxFlatNodes(db *sql.DB, appVer types.AppVersion, user_id interface{}) ([]models.FlatNode, error) {

	if appVer == types.V2_0 {
		return sqlite_app_v2_0.GetHomeFlatNodes(db, user_id)
	}
	if appVer == types.V3_3 {
		return sqlite_app_v3_3.GetInboxFlatNodes(db, user_id)
	}

	e := fmt.Errorf("GetInboxFlatNodes not implemented for app version %s\n", appVer)
	return nil, e
}

func GetUserNodes(db *sql.DB, appVer types.AppVersion, user *models.User) error {
	if appVer == types.V2_0 {
		return sqlite_app_v2_0.GetUserNodes(db, user)
	}
	if appVer == types.V3_3 {
		return sqlite_app_v3_3.GetUserNodes(db, user)
	}

	return fmt.Errorf("GetUserNodes not implemented for app version %s\n", appVer)
}
