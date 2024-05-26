package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
	"github.com/w3qxst1ck/cs2-grenades/internal/data"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		msxIdleConns int
		maxIdleTime  string
	}
	imagesDir string
	imagesUrl string
	limiter   struct {
		rps     float64 // request per seconds
		burst   int     // value burst per request
		enabled bool
	}
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
	wg     sync.WaitGroup
	cache  *cache.Cache
}

func main() {
	var cfg config

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	err := godotenv.Load()
	if err != nil {
		logger.Fatal(err)
	}

	// server
	apiPort, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		logger.Fatal(err)
	}

	flag.IntVar(&cfg.port, "port", apiPort, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// db configuration
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GRENADES_DB_DSN"), "PostreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgresSQL max open connections")
	flag.IntVar(&cfg.db.msxIdleConns, "db-max-idle-conns", 25, "PostgresSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgresSQL max connection idle time")

	// image
	flag.StringVar(&cfg.imagesDir, "images-directory", os.Getenv("API_IMAGES_DIR"), "Directory for saved images")
	flag.StringVar(&cfg.imagesUrl, "images-url", "/v1/image/", "Images url")

	// limiter
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 10, "Rate limiter maximum request per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	
	rateLimmiterEnable, err := strconv.ParseBool(os.Getenv("API_RATE_LIMITTER_ENABLE"))
	if err != nil {
		logger.Fatal(err)
	}
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", rateLimmiterEnable, "Enable rate limiter")

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	cache := cache.New(time.Minute*10, time.Minute*20)

	app := application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		cache:  cache,
	}

	err = app.serve()
	if err != nil {
		logger.Print(err)
	}
}
