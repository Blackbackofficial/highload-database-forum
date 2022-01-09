package repo

import (
	"context"
	"fmt"
	"forumI/internal/models"
	"forumI/internal/pkg/forum"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	SelectUserByNickname           = "select nickname, fullname, about, email from users where nickname=$1 limit 1;"
	SelectForumBySlug              = "select slug, \"user\", title, posts, threads from forum where slug=$1 limit 1;"
	InsertInForum                  = "insert into forum(slug, \"user\", title) values ($1, $2, $3);"
	InsertInThread                 = "insert into threads(title, author, created, forum, message, slug) values ($1, $2, $3, $4, $5, $6) returning *"
	SelectThreadSlug               = "select id, title, author, forum, message, votes, slug, created from threads where slug=$1 limit 1;"
	GetUsersOfForumDescNotNilSince = "select nickname, fullname, about, email from users_forum where slug=$1 and nickname < '%s' order by nickname desc limit nullif($2, 0)"
	GetUsersOfForumDescSinceNil    = "select nickname, fullname, about, email from users_forum where slug=$1 order by nickname desc limit nullif($2, 0)"
	GetUsersOfForumDescNil         = "select nickname, fullname, about, email from users_forum where slug=$1 and nickname > '%s' order by nickname limit nullif($2, 0)"
	GetThreadsSinceDescNotNil      = "select id, title, author, forum, message, votes, slug, created from threads where forum=$1 and created <= $2 order by created desc limit $3;"
	GetThreadsSinceDescNil         = "select id, title, author, forum, message, votes, slug, created from threads where forum=$1 and created >= $2 order by created asc limit $3;"
	GetThreadsDescNotNil           = "select id, title, author, forum, message, votes, slug, created from threads where forum=$1 order by created desc limit $2;"
	GetThreadsDescNil              = "select id, title, author, forum, message, votes, slug, created from threads where forum=$1 order by created asc limit $2;"
	SelectPostById                 = "select author, post, created_at, forum, isedited, parent, threads from posts where id = $1;"
	SelectThreadId                 = "select id, title, author, forum, message, votes, slug, created from thread where id=$1 LIMIT 1;"
)

type repoPostgres struct {
	Conn *pgxpool.Pool
}

func NewRepoPostgres(Conn *pgxpool.Pool) forum.Repository {
	return &repoPostgres{Conn: Conn}
}

func (r repoPostgres) GetUser(name string) (models.User, models.StatusCode) {
	var userM models.User
	row := r.Conn.QueryRow(context.Background(), SelectUserByNickname, name)
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
func (r repoPostgres) GetForum(slug string) (models.Forum, models.StatusCode) {
	forumM := models.Forum{}
	row := r.Conn.QueryRow(context.Background(), SelectForumBySlug, slug)
	err := row.Scan(&forumM.Slug, &forumM.User, &forumM.Title, &forumM.Posts, &forumM.Threads)
	if err != nil {
		return models.Forum{}, models.NotFound
	}
	return forumM, models.Okey
}

func (r repoPostgres) InThread(thread models.Thread) (models.Thread, error) {
	threadS := models.Thread{}
	row := r.Conn.QueryRow(context.Background(), InsertInThread, thread.Title,
		thread.Author, thread.Created, thread.Forum, thread.Message, thread.Slug)

	err := row.Scan(&threadS.Id, &threadS.Title, &threadS.Author, &threadS.Created,
		&threadS.Forum, &threadS.Message, &threadS.Slug, &threadS.Votes)
	if err != nil {
		return models.Thread{}, err
	}
	return threadS, nil
}

func (r repoPostgres) GetThreadSlug(slug string) (models.Thread, models.StatusCode) {
	threadS := models.Thread{}
	row := r.Conn.QueryRow(context.Background(), SelectThreadSlug, slug)
	err := row.Scan(&threadS.Id, &threadS.Title, &threadS.Author, &threadS.Forum,
		&threadS.Message, &threadS.Votes, &threadS.Slug, &threadS.Created)
	if err != nil {
		return models.Thread{}, models.NotFound
	}
	return threadS, models.Okey
}

func (r repoPostgres) GetUsersOfForum(forum models.Forum, limit string, since string, desc string) ([]models.User, models.StatusCode) {
	var query string
	if desc != "" {
		if since != "" {
			query = fmt.Sprintf(GetUsersOfForumDescNotNilSince, since)
		} else {
			query = GetUsersOfForumDescSinceNil
		}
	} else {
		query = fmt.Sprintf(GetUsersOfForumDescNil, since)
	}
	users := make([]models.User, 0)
	row, err := r.Conn.Query(context.Background(), query, forum.Slug, limit)

	if err != nil {
		return users, models.NotFound
	}

	defer row.Close() //??

	for row.Next() {
		user := models.User{}
		err = row.Scan(&user.NickName, &user.FullName, &user.About, &user.Email)
		if err != nil {
			return users, models.InternalError
		}
		users = append(users, user)
	}

	return users, models.Okey
}

func (r repoPostgres) GetThreadsOfForum(forum models.Forum, limit string, since string, desc string) ([]models.Thread, models.StatusCode) {
	var rows *pgx.Rows
	threads := make([]models.Thread, 0)

	if since != "" {
		if desc != "" {
			rows, err := r.Conn.Query(context.Background(), GetThreadsSinceDescNotNil, forum.Slug, since, limit)
			if err != nil {
				return threads, models.NotFound
			}
			defer rows.Close()
		} else {
			rows, err := r.Conn.Query(context.Background(), GetThreadsSinceDescNil, forum.Slug, since, limit)
			if err != nil {
				return threads, models.NotFound
			}
			defer rows.Close()
		}
	} else {
		if desc != "" {
			rows, err := r.Conn.Query(context.Background(), GetThreadsDescNotNil, forum.Slug, limit)
			if err != nil {
				return threads, models.NotFound
			}
			defer rows.Close()
		} else {
			rows, err := r.Conn.Query(context.Background(), GetThreadsDescNil, forum.Slug, limit)
			if err != nil {
				return threads, models.NotFound
			}
			defer rows.Close()
		}
	}

	for rows.Next() {
		threadS := models.Thread{}
		err := rows.Scan(&threadS.Id, &threadS.Title, &threadS.Author, &threadS.Forum, &threadS.Message,
			&threadS.Votes, &threadS.Slug, &threadS.Created)
		if err != nil {
			continue
		}
		threads = append(threads, threadS)
	}
	return threads, models.Okey
}

func (r repoPostgres) GetIdThread(id int) (models.Thread, models.StatusCode) {
	threadS := models.Thread{}
	row := r.Conn.QueryRow(context.Background(), SelectThreadId, id)

	err := row.Scan(&threadS.Id, &threadS.Title, &threadS.Author, &threadS.Forum, &threadS.Message,
		&threadS.Votes, &threadS.Slug, &threadS.Created)
	if err != nil {
		return models.Thread{}, models.NotFound
	}
	return threadS, models.Okey
}

func (r repoPostgres) GetFullPostInfo(posts models.PostFull, related []string) (models.PostFull, models.StatusCode) {
	post := models.Post{}
	postFull := models.PostFull{Author: nil, Forum: nil, Post: models.Post{}, Thread: nil}

	post.ID = posts.Post.ID

	row := r.Conn.QueryRow(context.Background(), SelectPostById, posts.Post.ID)
	err := row.Scan(&post.Author, &post.Message, &post.Created, &post.Forum, &post.IsEdited, &post.Parent, &post.Thread)

	if err != nil {
		return postFull, models.NotFound
	}

	postFull.Post = post

	for i := 0; i < len(related); i++ {
		if related[i] == "user" {
			user, _ := r.GetUser(post.Author)
			postFull.Author = &user
		}
		if related[i] == "forum" {

			forum, _ := r.GetForum(post.Forum)
			postFull.Forum = &forum

		}
		if related[i] == "thread" {
			thread, _ := r.GetIdThread(post.Thread)
			postFull.Thread = &thread

		}
	}
	return postFull, models.Okey
}
