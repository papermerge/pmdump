package database_app_v3_3

import (
	_ "github.com/mattn/go-sqlite3"

	models "github.com/papermerge/pmdump/models/app_v3_3"
	"github.com/papermerge/pmdump/types"
)

func GetPermissions(db *types.DBConn) ([]models.Permission, error) {
	rows, err := db.DB.Query("SELECT id, name, codename FROM permissions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []models.Permission

	for rows.Next() {
		var perm models.Permission

		err = rows.Scan(
			&perm.ID,
			&perm.Name,
			&perm.Codename,
		)
		if err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return perms, nil
}
