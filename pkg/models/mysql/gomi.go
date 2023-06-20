package mysql

import (
	"database/sql"

	"github.com/abhaysp95/gomibako/pkg/models"
)

type GomiModel struct {
	DB *sql.DB  // db pool dependency injection
}

// create new gomi
func (gm *GomiModel) Create(title, content, expires string) (int, error) {
	stmt := `insert into gomi(title, content, created, expires)
	values(?, ?, utc_timestamp(), date_add(utc_timestamp(), interval ? day))`

	res, err := gm.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// return the id of lastest created gomi
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}

	return int(id), nil
}

// get specific gomi
func (gm *GomiModel) Get(id int) (*models.Gomi, error) {
	return nil, nil
}

// 10 latest gomi
func (gm *GomiModel) Latest() ([]*models.Gomi, error) {
	return nil, nil
}