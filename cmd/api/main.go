package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/obynonwane/polygon_swiftlink_auth_api/data"
	"github.com/obynonwane/polygon_swiftlink_auth_api/token"
)

const webPort = "80"

var counts int64

type Config struct {
	Repo       data.Repository
	tokenMaker token.Maker
}

type TokenType struct {
	TokenSymmetricKey   string
	AccessTokenDuration time.Duration
}

func main() {

	//TokenSymmetricKey, AccessTokenDuration: Need to come from env
	pasetoDetail := &TokenType{
		TokenSymmetricKey:   os.Getenv("TOKEN_SYMETRIC_KEY"),
		AccessTokenDuration: 100 * time.Minute,
	}

	tokenMaker, err := token.NewPasetoMaker(pasetoDetail.TokenSymmetricKey)
	if err != nil {
		return
	}

	log.Println("Starting authentication service")

	//Connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to Postgres")
	}

	//setup config
	app := Config{
		Repo:       data.NewPostgresRepository(conn),
		tokenMaker: tokenMaker,
	}

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {

	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	sslmode := os.Getenv("DATABASE_SSLMODE")
	timezone := os.Getenv("DATABASE_TIMEZONE")
	connectTimeout := os.Getenv("DATABASE_CONNECT_TIMEOUT")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=%s connect_timeout=%s",
		host, port, user, password, dbname, sslmode, timezone, connectTimeout,
	)
	dsn := connStr
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Authentication Service Connected to Postgres ...")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("backing off for 2 seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}

func (app *Config) setupRepo(conn *sql.DB) {
	db := data.NewPostgresRepository(conn)
	app.Repo = db
}
