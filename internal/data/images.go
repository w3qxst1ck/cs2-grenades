package data

import (
	"context"
	"database/sql"
	"time"
)

type Image struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	ImageURL  string `json:"image_url"`
	GrenadeID int64  `json:"-"`
}

type ImageModel struct {
	DB *sql.DB
}

func (m ImageModel) Insert(image *Image) error {
	query := `
	INSERT INTO images (name, grenade_id)
	VALUES ($1, $2)
	RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, image.Name, image.GrenadeID).Scan(&image.ID)
}
