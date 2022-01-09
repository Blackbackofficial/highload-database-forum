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
	conn := "postgres://postgres:password@127.0.0.1:5432/bd?sslmode=disable"
	//conn := "postgres://docker:docker@127.0.0.1:5432/docker?sslmode=disable"
	pool, err := pgxpool.Connect(context.Background(), conn)
	if err != nil {
		log.Fatal("No connection to postgres", err)
	}

	fRepo := repo.NewRepoPostgres(pool)
	fUsecase := usecase.NewRepoUsecase(fRepo)
	fHandler := delivery.NewForumHandler(fUsecase)

	forum := muxRoute.PathPrefix("/api").Subrouter()
	{
		forum.HandleFunc("/forum/create", fHandler.CreateForum).Methods(http.MethodPost)
		forum.HandleFunc("/forum/{slug}/details", fHandler.ForumInfo).Methods(http.MethodGet)
		forum.HandleFunc("/forum/{slug}/create", fHandler.CreateThreadsForum).Methods(http.MethodPost)
		forum.HandleFunc("/forum/{slug}/users", fHandler.GetUsersForum).Methods(http.MethodGet)
		forum.HandleFunc("/forum/{slug}/threads", fHandler.GetThreadsForum).Methods(http.MethodGet)

		forum.HandleFunc("/post/{id}/details", fHandler.GetPostInfo).Methods(http.MethodGet)
		forum.HandleFunc("/post/{id}/details", fHandler.UpdatePostInfo).Methods(http.MethodGet)

		forum.HandleFunc("/service/clear", fHandler.GetClear).Methods(http.MethodPost)
		forum.HandleFunc("/service/status", fHandler.GetStatus).Methods(http.MethodGet)
	}

	http.Handle("/", muxRoute)
	log.Print(http.ListenAndServe(":5000", muxRoute))
}
