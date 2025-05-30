package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Expires time.Time `json:"expires"`
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	stmt := `
		INSERT INTO snippets (title, content, expires)
		VALUES ($1, $2, CURRENT_TIMESTAMP + make_interval(days => $3))
		RETURNING id
	`
	var lastInsertId int
	err := m.DB.QueryRow(stmt, title, content, expires).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `
	SELECT id, title, content, expires 
	FROM snippets
	WHERE expires > CURRENT_TIMESTAMP and id = $1
	`
	var s Snippet
	err := m.DB.QueryRow(stmt, id).Scan(
		&s.ID,
		&s.Title,
		&s.Content,
		&s.Expires,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `
	SELECT id, title, content, expires FROM snippets
	WHERE expires > CURRENT_TIMESTAMP ORDER BY id DESC LIMIT 10
	`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Expires)
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
