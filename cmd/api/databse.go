package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	// устанавливаем максимально возможное кол-во соединений (in-use + idle) с БД, default - no limit
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	// устанавливаем максимально возможное кол-во idle соединений, default - no limit
	db.SetMaxIdleConns(cfg.db.msxIdleConns)

	// парсим "15m" в тип time.Duration с помощью time.ParseDuration и уст. макс. idle timeout
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
