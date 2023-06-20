package mysql

import (
	"database/sql"

	"github.com/abhaysp95/gomibako/pkg/models"
)

type GomiModel struct {
	DB *sql.DB // db pool dependency injection
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
	stmt := `select id, title, content, created, expires from gomi
	where expires > utc_timestamp() and id = ?`

	g := &models.Gomi{}

	// return atmost on matching row
	if err := gm.DB.QueryRow(stmt, id).Scan(&g.Id, &g.Title, &g.Content, &g.Created, &g.Expires); err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return g, nil
}

// 10 latest gomi
func (gm *GomiModel) Latest() ([]*models.Gomi, error) {
	stmt := `select id, title, content, created, expires from gomi
	where expires > utc_timestamp() order by created desc limit 10`

	var gs []*models.Gomi

	rows, err := gm.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// CRUCIAL: if resultset is open, db connection can't be closed and thus
	// pool will keep using the resource
	defer rows.Close()

	for rows.Next() {
		g := &models.Gomi{}
		if err = rows.Scan(&g.Id, &g.Title, &g.Content, &g.Created, &g.Expires); err != nil {
			if err == sql.ErrNoRows {
				return nil, models.ErrNoRecord
			} else {
				return nil, err
			}
		}

		gs = append(gs, g)
	}

	return gs, nil
}
