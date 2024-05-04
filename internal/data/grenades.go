package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Grenade struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Side        string `json:"side"`
	Version     int32  `json:"version"`
}

type GrenadeModel struct {
	DB *sql.DB
}


func(m GrenadeModel) Get(id int64) (*Grenade, error) {
	query := `
	SELECT id, title, description, type, side, version
	FROM grenades 
	WHERE id = $1`

	var grenade Grenade

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&grenade.ID,
		&grenade.Title,
		&grenade.Description,
		&grenade.Type,
		&grenade.Side,
		&grenade.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &grenade, nil
	
}
