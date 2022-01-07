package repo

import (
	"context"
	"forumI/internal/models"
	"forumI/internal/pkg/forum"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	SelectUserByNickname = "select nickname, fullname, about, email from users where nickname=$1 limit 1;"
	SelectForumBySlug    = "select slug, \"user\", title, posts, threads From forum where slug=$1 limit 1;"
	InsertInForum        = "insert into forum(slug, \"user\", title) values ($1, $2, $3);"
)

type repoPostgres struct {
	Conn *pgxpool.Pool
}

func NewRepoPostgres(Conn *pgxpool.Pool) forum.Repository {
	return &repoPostgres{Conn: Conn}
}

func (r repoPostgres) GetUser(forum models.Forum) (models.User, models.StatusCode) {
	var userM models.User
	row := r.Conn.QueryRow(context.Background(), SelectUserByNickname, forum.User)
	err := row.Scan(&userM.NickName, &userM.FullName, &userM.About, &userM.Email)
	if err != nil {
		return models.User{}, models.NotFound
	}
	return userM, models.Okey
}

func (r repoPostgres) InForum(forum models.Forum) error {
	_, err := r.Conn.Exec(context.Background(), InsertInForum, forum.Slug, forum.User, forum.Title)
	if err != nil {
		return err
	}
	return nil
}
func (r repoPostgres) GetForum(forum models.Forum) (models.Forum, models.StatusCode) {
	var forumM models.Forum
	row := r.Conn.QueryRow(context.Background(), SelectForumBySlug, forum.Slug)
	err := row.Scan(&forum.Slug, &forum.User, &forum.Title, &forum.Posts, &forum.Threads)
	if err != nil {
		return models.Forum{}, models.NotFound
	}
	return forumM, models.Okey
}
