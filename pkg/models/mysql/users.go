package mysql

import (
	"database/sql"

	"github.com/abhaysp95/gomibako/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB  *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, passwd, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stmt, name, email, string(hash))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {  // Duplicate entry error code
				return models.ErrDuplicateEmail
			}
		}
	}

	return err
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	stmt := `select id, passwd from users where email = ?`
	rows := m.DB.QueryRow(stmt, email)

	var id int
	var hashedPasswd string
	err := rows.Scan(&id, &hashedPasswd)
	if err != nil {
		return 0, models.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
