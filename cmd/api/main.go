package main

import (
	"backend/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Version of application
const version = "1.0.0"

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type application struct {
	config config
	logger *log.Logger
	models models.Models
}

func main() {
	var cfg config

	// flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	// flag.StringVar(&cfg.env, "env", "development", "Application enviornment (development|production)")
	// flag.StringVar(&cfg.db.dsn, "dsn", "postgres://toshi:postoshiword@postgres/go-movies?sslmode=disable", "Postgres connection string")
	// flag.Parse()

	err := setConfig(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	log.Println("Connecting to database...")
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	log.Println("Connected to database!")

	// db, err := openDB(cfg)
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	// // Once you open a connection pool, you need to defer closing it.
	// defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server on port:", cfg.port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func setConfig(cfg *config) error {
	//----------------------------------------------------------------
	// Read .env parameters using gotodoenv library
	//----------------------------------------------------------------
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("cannot read .env file")
		return err
	}

	cfg.port = os.Getenv("PORT")
	cfg.env = os.Getenv("ENV")

	dbHost := os.Getenv("POSTGRESQL_HOST")
	dbName := os.Getenv("POSTGRESQL_DBNAME")
	dbUser := os.Getenv("POSTGRESQL_USERNAME")
	dbPass := os.Getenv("POSTGRESQL_PASSWORD")
	dbPort := os.Getenv("POSTGRESQL_PORT")
	dbSSL := os.Getenv("POSTGRESQL_SSL")

	cfg.db.dsn = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", dbHost, dbPort, dbName, dbUser, dbPass, dbSSL)

	return nil
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
