package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Snippet struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Expires   time.Time `json:"expires"`
}

type SnippetModel struct {
	Pool *pgxpool.Pool
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	stmt := `
		INSERT INTO snippets (title, content, expires)
		VALUES ($1, $2, CURRENT_TIMESTAMP + make_interval(days => $3))
		RETURNING id
	`
	var lastInsertId int
	err := m.Pool.QueryRow(context.Background(), stmt, title, content, expires).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `
	SELECT id, title, content, created_at, expires 
	FROM snippets
	WHERE expires > CURRENT_TIMESTAMP and id = $1
	`
	var s Snippet
	err := m.Pool.QueryRow(context.Background(), stmt, id).Scan(
		&s.ID,
		&s.Title,
		&s.Content,
		&s.CreatedAt,
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
	SELECT id, title, content, created_at, expires FROM snippets
	WHERE expires > CURRENT_TIMESTAMP ORDER BY id DESC LIMIT 10
	`
	rows, err := m.Pool.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.CreatedAt, &s.Expires)
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
