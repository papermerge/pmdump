package sqlite_db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	sqlite_app_v2_0 "github.com/papermerge/pmdump/database/sqlite/app_v2_0"
	sqlite_app_v3_3 "github.com/papermerge/pmdump/database/sqlite/app_v3_3"
	"github.com/papermerge/pmdump/models"
	"github.com/papermerge/pmdump/types"
	"github.com/papermerge/pmdump/utils"
)

func Open(dburl string) (*sql.DB, error) {
	/* at this point `sql.Open` won't complain if dburl
	   is a path to folder which will result in confusing error
	   message. Double check now that dburl points to a file
	*/
	if !utils.IsReadableFile(dburl) {
		return nil, fmt.Errorf("%q is not a readable file", dburl)
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

	e := fmt.Errorf("GetUsers not implemented for app version %q", appVer)
	return nil, e
}

func GetHomeFlatNodes(db *sql.DB, appVer types.AppVersion, user_id interface{}) ([]models.FlatNode, error) {

	if appVer == types.V2_0 {
		return sqlite_app_v2_0.GetHomeFlatNodes(db, user_id)
	}
	if appVer == types.V3_3 {
		return sqlite_app_v3_3.GetHomeFlatNodes(db, user_id)
	}

	e := fmt.Errorf("GetHomeFlatNodes not implemented for app version %q", appVer)
	return nil, e
}

func GetInboxFlatNodes(db *sql.DB, appVer types.AppVersion, user_id interface{}) ([]models.FlatNode, error) {

	if appVer == types.V2_0 {
		return sqlite_app_v2_0.GetHomeFlatNodes(db, user_id)
	}
	if appVer == types.V3_3 {
		return sqlite_app_v3_3.GetInboxFlatNodes(db, user_id)
	}

	e := fmt.Errorf("GetInboxFlatNodes not implemented for app version %q", appVer)
	return nil, e
}

func GetUserNodes(db *sql.DB, appVer types.AppVersion, user *models.User) error {
	if appVer == types.V2_0 {
		return sqlite_app_v2_0.GetUserNodes(db, user)
	}
	if appVer == types.V3_3 {
		return sqlite_app_v3_3.GetUserNodes(db, user)
	}

	return fmt.Errorf("GetUserNodes not implemented for app version %q", appVer)
}

func GetDocumentPageRows(db *sql.DB, appVer types.AppVersion, user_id interface{}) ([]models.DocumentPageRow, error) {
	if appVer == types.V2_0 {
		return sqlite_app_v2_0.GetDocumentPageRows(db, user_id)
	}
	if appVer == types.V3_3 {
		return sqlite_app_v3_3.GetDocumentPageRows(db, user_id)
	}

	return nil, fmt.Errorf("GetDocumentPageRows not implemented for app version %q", appVer)
}
