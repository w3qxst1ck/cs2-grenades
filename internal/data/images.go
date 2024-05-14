package data

import (
	"context"
	"database/sql"
	"errors"
	"mime/multipart"
	"strings"
	"time"

	"github.com/w3qxst1ck/cs2-grenades/internal/validator"
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

func VlidateImage(fileHeader *multipart.FileHeader, v *validator.Validator) {
	v.Check(fileHeader.Size < 20_000_000, "grenadeImage_size", "file size must be less than 20MB")
	v.Check(v.In(strings.Split(fileHeader.Filename, ".")[1], []string{"jpg", "jpeg", "png"}), "grenadeImage_extension", "file extension must be jpeg|jpg|png")
}

func (m ImageModel) Get(id int64) (*Image, error) {
	query := `
	SELECT id, name FROM images
	WHERE id = $1`

	var image Image

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&image.ID,
		&image.Name,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &image, nil
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

func (m ImageModel) GetByGrenadeID(grenadeId int64) ([]*Image, error) {
	query := `
	SELECT id, name FROM images
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

		err := rows.Scan(
			&image.ID,
			&image.Name,
		)
		if err != nil {
			return nil, err
		}
		image.GrenadeID = grenadeId

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

func (m ImageModel) Delete(id int64) error {
	query := `
	DELETE FROM images
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
