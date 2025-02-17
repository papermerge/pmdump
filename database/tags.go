package database

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/papermerge/pmg-dump/models"
)

func GetTags(db *sql.DB) ([]models.Tag, error) {
	rows, err := db.Query("SELECT id, name, pinned, description, fg_color, bg_color FROM core_tag")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag

	for rows.Next() {
		var tag models.Tag
		err = rows.Scan(&tag.ID, &tag.Name, &tag.Pinned, &tag.Description, &tag.FGColor, &tag.BGColor)
		if err != nil {
			return nil, err
		}
		tag.UUID = uuid.New()
		tags = append(tags, tag)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
