package main

import (
	"context"
	"forumI/internal/pkg/forum/delivery"
	"forumI/internal/pkg/forum/repo"
	"forumI/internal/pkg/forum/usecase"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
)

func main() {
	muxRoute := mux.NewRouter()
	//conn := "postgres://postgres:password@127.0.0.1:5432/db?sslmode=disable"
	conn := "postgres://docker:docker@127.0.0.1:5432/docker?sslmode=disable"
	pool, err := pgxpool.Connect(context.Background(), conn)
	if err != nil {
		log.Fatal(pool)
	}

	fRepo := repo.NewRepoPostgres(pool)
	fUsecase := usecase.NewRepoUsecase(fRepo)
	fHandler := delivery.NewForumHandler(fUsecase)

	forum := muxRoute.PathPrefix("/api").Subrouter()
	{
		forum.HandleFunc("/forum/create", fHandler.CreateForum).Methods(http.MethodPost)
	}

	http.Handle("/", muxRoute)
	log.Print(http.ListenAndServe(":5000", muxRoute))

}
