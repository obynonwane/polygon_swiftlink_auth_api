package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/obynonwane/polygon_swiftlink_auth_api/data"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	app := Config{}

	log.Printf("starting authentication service on port %s\n", webPort)
	//define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	//start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
