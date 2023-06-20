package mysql

import (
	"database/sql"

	"github.com/abhaysp95/gomibako/pkg/models"
)

type GomiModel struct {
	DB *sql.DB
}

// create new gomi
func (gm *GomiModel) Create(title, content, expires string) (int, error) {
	return 0, nil
}

// get specific gomi
func (gm *GomiModel) Get(id int) (*models.Gomi, error) {
	return nil, nil
}

// 10 latest gomi
func (gm *GomiModel) Latest() ([]*models.Gomi, error) {
	return nil, nil
}
