package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	// кастомная ошибка которую возвращаем когда конфликт во внесении изменений в БД
	ErrEditConflict = errors.New("edit conflict")
)

type Models struct {
	Grenades GrenadeModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Grenades: GrenadeModel{DB: db},
	}
}