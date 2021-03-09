package mysql

import (
	"database/sql"
	"errors"

	"github.com/zappel/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := "INSERT INTO snippets(title, content, created, expires) VALUES(?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY))"

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {

	//cara 1 untuk query atau bisa cara yang bawah
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id =?"
	row := m.DB.QueryRow(stmt, id)
	/*s := &models.Snippet{}

	err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	err := m.DB.QueryRow("SELECT ...", id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	*/

	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// hal 214 “In the code above we’ve used the errors.Is() function — which was introduced in Go 1.13 — to check wether an error matches a specific value (in our case, we wanted to check if an error matched sql.ErrNoRows).”
			//error.IS (keknya untuk bandingin apakah error = error sql) ini function dari golang untuk handle error dari sql?
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	 WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires) //disini jadi dari hasil scan db nya di copy ke model s yang di buat di package models
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
