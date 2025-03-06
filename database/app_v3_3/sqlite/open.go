package sqlite_app_v3_3

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
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
