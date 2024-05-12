package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/w3qxst1ck/cs2-grenades/internal/validator"
)

type Grenade struct {
	ID          int64   `json:"id"`
	Map         string  `json:"map"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Side        string  `json:"side"`
	Version     int32   `json:"version"`
	Images      []Image `json:"images,omitempty"`
}

type GrenadeModel struct {
	DB *sql.DB
}

func ValidateGrenade(grenade *Grenade, v *validator.Validator) {
	v.Check(grenade.Map != "", "map", "must be provided")
	v.Check(len(grenade.Map) <= 100, "map", "must not be grater than 100 bytes")

	v.Check(grenade.Title != "", "title", "must be provided")
	v.Check(len(grenade.Title) <= 500, "title", "must not be grater than 500 bytes")

	v.Check(len(grenade.Description) <= 700, "description", "must not be grater than 700 bytes")

	v.Check(grenade.Type != "", "type", "must be provided")
	v.Check(v.In(grenade.Type, []string{"smoke", "molotov", "he", "flash", "decoy"}), "type", "value of type must be smoke|molotov|he|flash|decoy")

	v.Check(grenade.Side != "", "side", "must be provided")
	v.Check(v.In(grenade.Side, []string{"CT", "T"}), "side", "value of side must be T or CT")
}

func (m GrenadeModel) Get(id int64) (*Grenade, error) {
	query := `
	SELECT g.id, g.map, g.title, g.description, g.type, g.side, g.version, string_agg(i.name, ',')
	FROM grenades g
	LEFT JOIN images i
	ON g.id = i.grenade_id
	WHERE g.id = $1
	GROUP BY i.grenade_id, g.id`

	var grenade Grenade
	var imagesNames string

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&grenade.ID,
		&grenade.Map,
		&grenade.Title,
		&grenade.Description,
		&grenade.Type,
		&grenade.Side,
		&grenade.Version,
		&imagesNames,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	for _, name := range strings.Split(imagesNames, ",") {
		grenade.Images = append(grenade.Images, Image{
			Name: name,
		})
	}

	return &grenade, nil
}

func (m GrenadeModel) Insert(grenade *Grenade) error {
	query := `
	INSERT INTO grenades (map, title, description, type, side)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{grenade.Map, grenade.Title, grenade.Description, grenade.Type, grenade.Side}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&grenade.ID, &grenade.Version)
}

func (m GrenadeModel) Update(grenade *Grenade) error {
	query := `
	UPDATE grenades 
	SET map=$1, title=$2, description=$3, type=$4, side=$5, version=version + 1
	WHERE id=$6 AND version=$7
	RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{
		grenade.Map,
		grenade.Title,
		grenade.Description,
		grenade.Type,
		grenade.Side,
		grenade.ID,
		grenade.Version,
	}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&grenade.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m GrenadeModel) Delete(id int64) error {
	query := `
	DELETE FROM grenades
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

func (m GrenadeModel) GetAll(csMap string, side string, grenType string, filters Filters) ([]*Grenade, error) {
	query := fmt.Sprintf(`
	SELECT g.id, g.map, g.title, g.description, g.type, g.side, g.version, string_agg(i.name, ',')
	FROM grenades g
	LEFT JOIN images i
	ON g.id = i.grenade_id
	WHERE (map = $1 OR $1 = '') AND (side = $2 OR $2 = '') AND (type = $3 OR $3 = '')
	GROUP BY i.grenade_id, g.id
	ORDER BY %s %s, id ASC`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, csMap, side, grenType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grenades := []*Grenade{}

	for rows.Next() {
		var grenade Grenade
		var imagesNames string

		err := rows.Scan(
			&grenade.ID,
			&grenade.Map,
			&grenade.Title,
			&grenade.Description,
			&grenade.Type,
			&grenade.Side,
			&grenade.Version,
			&imagesNames,
		)
		if err != nil {
			return nil, err
		}

		for _, name := range strings.Split(imagesNames, ",") {
			grenade.Images = append(grenade.Images, Image{
				Name: name,
			})
		}

		grenades = append(grenades, &grenade)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return grenades, nil
}
