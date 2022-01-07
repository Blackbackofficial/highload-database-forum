package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
)

func main() {
	conn := "postgres://postgres:password@127.0.0.1:5432/db?sslmode=disable"
	//conn := "postgres://docker:docker@127.0.0.1:5432/docker?sslmode=disable"
	pool, err := pgxpool.Connect(context.Background(), conn)
	if err != nil {
		log.Fatal(pool)
	}

	route := mux.NewRouter()
	fmt.Println(route)

	// ручки

	// lsining port
	http.Handle("/", route)
	log.Print(http.ListenAndServe(":5000", route))

}
