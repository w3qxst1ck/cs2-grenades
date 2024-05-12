package data

import (
	"context"
	"database/sql"
	"time"
)

type Image struct {
	ID        int64  `json:"id,omitempty"`
	Name      string `json:"name"`
	GrenadeID int64  `json:"-"`
	ImageURL  string `json:"image_url"`
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

func (m ImageModel) Get(grenadeId int64) ([]*Image, error) {
	query := `
	SELECT name FROM images
	WHERE grenade_id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	images := []*Image{}

	rows, err := m.DB.QueryContext(ctx, query, grenadeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var image Image

		err := rows.Scan(&image.Name)
		if err != nil {
			return nil, err
		}

		images = append(images, &image)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

func (m ImageModel) GetAll() ([]*Image, error) {
	query := `
	SELECT name, grenade_id
	FROM images`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := []*Image{}

	for rows.Next() {
		var image Image

		err := rows.Scan(
			&image.Name,
			&image.GrenadeID,
		)
		if err != nil {
			return nil, err
		}

		images = append(images, &image)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}