package repo

import (
	"context"
	"fmt"
	"forumI/internal/models"
	"forumI/internal/pkg/forum"
	"forumI/internal/pkg/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
	"time"
)

const (
	SelectUserByNickname            = "select nickname, fullname, about, email from users where nickname=$1 limit 1;"
	SelectForumBySlug               = "select slug, \"user\", title, posts, threads from forum where slug=$1 limit 1;"
	InsertInForum                   = "insert into forum(slug, \"user\", title) values ($1, $2, $3);"
	InsertInThread                  = "insert into threads(title, author, created, forum, message, slug) values ($1, $2, $3, $4, $5, $6) returning *"
	SelectThreadSlug                = "select id, title, author, forum, message, votes, slug, created from threads where slug=$1 limit 1;"
	GetUsersOfForumDescNotNilSince  = "select nickname, fullname, about, email from users_forum where slug=$1 and nickname < '%s' order by nickname desc limit nullif($2, 0)"
	GetUsersOfForumDescSinceNil     = "select nickname, fullname, about, email from users_forum where slug=$1 order by nickname desc limit nullif($2, 0)"
	GetUsersOfForumDescNil          = "select nickname, fullname, about, email from users_forum where slug=$1 and nickname > '%s' order by nickname limit nullif($2, 0)"
	GetThreadsSinceDescNotNil       = "select id, title, author, forum, message, votes, slug, created from threads where forum=$1 and created <= $2 order by created desc limit $3;"
	GetThreadsSinceDescNil          = "select id, title, author, forum, message, votes, slug, created from threads where forum=$1 and created >= $2 order by created asc limit $3;"
	GetThreadsDescNotNil            = "select id, title, author, forum, message, votes, slug, created from threads where forum=$1 order by created desc limit $2;"
	GetThreadsDescNil               = "select id, title, author, forum, message, votes, slug, created from threads where forum=$1 order by created asc limit $2;"
	SelectPostById                  = "select author, post, created_at, forum, isedited, parent, threads from posts where id = $1;"
	SelectThreadId                  = "select id, title, author, forum, message, votes, slug, created from threads where id=$1 LIMIT 1;"
	UpdatePostMessage               = "update posts set message=coalesce(nullif($1, ''), message), isedited = case when $1 = '' or message = $1 then isedited else true end where id=$2 returning *"
	ClearAll                        = "truncate table users, forum, threads, posts, votes, users_forum CASCADE;"
	SelectCountUsers                = "select count(*) from users;"
	SelectCountForum                = "select count(*) from forum;"
	SelectCountThreads              = "select count(*) from threads;"
	SelectCountPosts                = "select count(*) from posts;"
	InsertManyPosts                 = "insert into posts(author, created, forum, message, parent, thread) values"
	UpdateThread                    = "update threads set title=coalesce(nullif($1, ''), title), message=coalesce(nullif($2, ''), message) where %s returning *"
	SelectPostSinceDescNotNil       = "select id, author, created, forum, isedited, message, parent, thread from posts where thread=$1 order by id desc limit $2;"
	SelectPostSinceDescNil          = "select id, author, created, forum, isedited, message, parent, thread from posts where thread=$1 order by id limit $2;"
	SelectPostDescNotNil            = "select id, author, created, forum, isedited, message, parent, thread from posts where thread=$1 and id < $2 order by id desc limit $3;"
	SelectPostDescNil               = "select id, author, created, forum, isedited, message, parent, thread from posts where thread=$1 and id > $2 order by id limit $3;"
	SelectPostTreeSinceDescNotNil   = "select id, author, created, forum, isedited, message, parent, thread from posts where thread=$1 order by path desc, id desc limit $2;"
	SelectPostTreeSinceDescNil      = "select id, author, created, forum, isedited, message, parent, thread from posts where thread=$1 order by path asc, id asc limit $2;"
	SelectPostTreeDescNotNil        = "select id, author, created, forum, isedited, message, parent, thread from posts where thread=$1 and path < (select path from posts where id = $2) order by path desc, id desc limit $3;"
	SelectPostTreeDescNil           = "select id, author, created, forum, isedited, message, parent, thread from posts where thread=$1 and path > (select path from posts where id = $2) order by path asc, id asc limit $3;"
	SelectPostParentSinceDescNotNil = "select id, author, created, forum, isedited, message, parent, thread from posts where path[1] in (select id from posts where thread = $1 and parent is null order by id desc limit $2) order by path[1] desc, path, id;"
	SelectPostParentSinceDescNil    = "select id, author, created, forum, isedited, message, parent, thread from posts where path[1] in (select id from posts where thread = $1 and parent is null order by id limit $2) order by path, id;"
	SelectPostParentDescNotNil      = "select id, author, created, forum, isedited, message, parent, thread from posts where path[1] IN (select id from posts where thread = $1 and parent is null and path[1] < (select path[1] from posts where id = $2) order by id desc limit $3) order by path[1] desc, path, id;"
	SelectPostParentDescNil         = "select id, author, created, forum, isedited, message, parent, thread from posts where path[1] IN (select id from posts where thread = $1 and parent is null and path[1] > (select path[1] from posts where id = $2) order by id asc limit $3) order by path, id;"
	UpdateVote                      = "update votes set voice=$1 where author=$2 and thread=$3;"
	InsertVote                      = "insert into votes(author, voice, thread) values ($1, $2, $3);"
)

type repoPostgres struct {
	Conn *pgxpool.Pool
}

func NewRepoPostgres(Conn *pgxpool.Pool) forum.Repository {
	return &repoPostgres{Conn: Conn}
}

func (r *repoPostgres) GetUser(name string) (models.User, models.StatusCode) {
	var userM models.User
	row := r.Conn.QueryRow(context.Background(), SelectUserByNickname, name)
	err := row.Scan(&userM.NickName, &userM.FullName, &userM.About, &userM.Email)
	if err != nil {
		return models.User{}, models.NotFound
	}
	return userM, models.Okey
}

func (r *repoPostgres) InForum(forum models.Forum) error {
	_, err := r.Conn.Exec(context.Background(), InsertInForum, forum.Slug, forum.User, forum.Title)
	if err != nil {
		return err
	}
	return nil
}

func (r *repoPostgres) GetForum(slug string) (models.Forum, models.StatusCode) {
	forumM := models.Forum{}
	row := r.Conn.QueryRow(context.Background(), SelectForumBySlug, slug)
	err := row.Scan(&forumM.Slug, &forumM.User, &forumM.Title, &forumM.Posts, &forumM.Threads)
	if err != nil {
		return models.Forum{}, models.NotFound
	}
	return forumM, models.Okey
}

func (r *repoPostgres) InThread(thread models.Thread) (models.Thread, error) {
	threadS := models.Thread{}
	row := r.Conn.QueryRow(context.Background(), InsertInThread, thread.Title,
		thread.Author, thread.Created, thread.Forum, thread.Message, thread.Slug)

	err := row.Scan(&threadS.ID, &threadS.Title, &threadS.Author, &threadS.Created,
		&threadS.Forum, &threadS.Message, &threadS.Slug, &threadS.Votes)
	if err != nil {
		return models.Thread{}, err
	}
	return threadS, nil
}

func (r *repoPostgres) GetThreadSlug(slug string) (models.Thread, models.StatusCode) {
	threadS := models.Thread{}
	row := r.Conn.QueryRow(context.Background(), SelectThreadSlug, slug)
	err := row.Scan(&threadS.ID, &threadS.Title, &threadS.Author, &threadS.Forum,
		&threadS.Message, &threadS.Votes, &threadS.Slug, &threadS.Created)
	if err != nil {
		return models.Thread{}, models.NotFound
	}
	return threadS, models.Okey
}

func (r *repoPostgres) GetUsersOfForum(forum models.Forum, limit string, since string, desc string) ([]models.User, models.StatusCode) {
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

func (r *repoPostgres) GetThreadsOfForum(forum models.Forum, limit string, since string, desc string) ([]models.Thread, models.StatusCode) {
	threads := make([]models.Thread, 0)

	if since != "" {
		if desc != "" {
			rows, err := r.Conn.Query(context.Background(), GetThreadsSinceDescNotNil, forum.Slug, since, limit)
			if err != nil {
				return threads, models.NotFound
			}
			defer rows.Close()
			for rows.Next() {
				threadS := models.Thread{}
				err := rows.Scan(&threadS.ID, &threadS.Title, &threadS.Author, &threadS.Forum, &threadS.Message,
					&threadS.Votes, &threadS.Slug, &threadS.Created)
				if err != nil {
					continue
				}
				threads = append(threads, threadS)
			}
		} else {
			rows, err := r.Conn.Query(context.Background(), GetThreadsSinceDescNil, forum.Slug, since, limit)
			if err != nil {
				return threads, models.NotFound
			}
			defer rows.Close()
			for rows.Next() {
				threadS := models.Thread{}
				err := rows.Scan(&threadS.ID, &threadS.Title, &threadS.Author, &threadS.Forum, &threadS.Message,
					&threadS.Votes, &threadS.Slug, &threadS.Created)
				if err != nil {
					continue
				}
				threads = append(threads, threadS)
			}
		}
	} else {
		if desc != "" {
			rows, err := r.Conn.Query(context.Background(), GetThreadsDescNotNil, forum.Slug, limit)
			if err != nil {
				return threads, models.NotFound
			}
			defer rows.Close()
			for rows.Next() {
				threadS := models.Thread{}
				err := rows.Scan(&threadS.ID, &threadS.Title, &threadS.Author, &threadS.Forum, &threadS.Message,
					&threadS.Votes, &threadS.Slug, &threadS.Created)
				if err != nil {
					continue
				}
				threads = append(threads, threadS)
			}
		} else {
			rows, err := r.Conn.Query(context.Background(), GetThreadsDescNil, forum.Slug, limit)
			if err != nil {
				return threads, models.NotFound
			}
			defer rows.Close()
			for rows.Next() {
				threadS := models.Thread{}
				err := rows.Scan(&threadS.ID, &threadS.Title, &threadS.Author, &threadS.Forum, &threadS.Message,
					&threadS.Votes, &threadS.Slug, &threadS.Created)
				if err != nil {
					continue
				}
				threads = append(threads, threadS)
			}
		}
	}
	return threads, models.Okey
}

func (r *repoPostgres) GetIdThread(id int) (models.Thread, models.StatusCode) {
	threadS := models.Thread{}
	row := r.Conn.QueryRow(context.Background(), SelectThreadId, id)

	err := row.Scan(&threadS.ID, &threadS.Title, &threadS.Author, &threadS.Forum,
		&threadS.Message, &threadS.Votes, &threadS.Slug, &threadS.Created)
	if err != nil {
		return models.Thread{}, models.NotFound
	}
	return threadS, models.Okey
}

func (r *repoPostgres) GetFullPostInfo(posts models.PostFull, related []string) (models.PostFull, models.StatusCode) {
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
		if "user" == related[i] {
			user, _ := r.GetUser(post.Author)
			postFull.Author = &user
		}
		if "forum" == related[i] {
			forumS, _ := r.GetForum(post.Forum)
			postFull.Forum = &forumS
		}
		if "thread" == related[i] {
			thread, _ := r.GetIdThread(post.Thread)
			postFull.Thread = &thread

		}
	}
	return postFull, models.Okey
}

func (r *repoPostgres) UpdatePostInfo(postOne models.Post, postUpdate models.PostUpdate) (models.Post, models.StatusCode) {
	row := r.Conn.QueryRow(context.Background(), UpdatePostMessage, postUpdate.Message, postOne.ID)
	err := row.Scan(&postOne.ID, &postOne.Author, &postOne.Created, &postOne.Forum,
		&postOne.IsEdited, &postOne.Message, &postOne.Parent, &postOne.Thread, &postOne.Path)
	if err != nil {
		return postOne, models.NotFound
	}
	return postOne, models.Okey
}

func (r *repoPostgres) GetClear() models.StatusCode {
	_, err := r.Conn.Exec(context.Background(), ClearAll)
	if err != nil {
		return models.InternalError
	}
	return models.Okey
}

func (r *repoPostgres) GetStatus() models.Status {
	statusS := models.Status{}
	row := r.Conn.QueryRow(context.Background(), SelectCountUsers)
	err := row.Scan(&statusS.User)
	if err != nil {
		statusS.User = 0
	}

	row = r.Conn.QueryRow(context.Background(), SelectCountForum)
	err = row.Scan(&statusS.Forum)
	if err != nil {
		statusS.Forum = 0
	}

	row = r.Conn.QueryRow(context.Background(), SelectCountThreads)
	err = row.Scan(&statusS.Thread)
	if err != nil {
		statusS.Thread = 0
	}

	row = r.Conn.QueryRow(context.Background(), SelectCountPosts)
	err = row.Scan(&statusS.Post)
	if err != nil {
		statusS.Post = 0
	}
	return statusS
}

func (r *repoPostgres) InPosts(postsS []models.Post, thread models.Thread) ([]models.Post, error) {
	rowQuery := InsertManyPosts
	data := make([]interface{}, 0)
	createdTime := time.Now()
	for i, onePost := range postsS {
		values := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		rowQuery += values
		data = append(data, thread.ID)
		data = append(data, onePost.Parent)
		data = append(data, onePost.Author)
		data = append(data, onePost.Message)
		data = append(data, thread.Forum)
		data = append(data, createdTime)
	}

	rowQuery = strings.TrimSuffix(rowQuery, ",")
	rowQuery += ` RETURNING id, isEdited, forum, thread, created;`

	rows, err := r.Conn.Query(context.Background(), rowQuery, data...)
	if err != nil {
		return nil, err
	}

	for i := range postsS {
		if rows.Next() {
			err := rows.Scan(&postsS[i].ID, &postsS[i].IsEdited, &postsS[i].Forum, &postsS[i].Thread, &postsS[i].Created)
			if err != nil {
				return nil, utils.Conflict
			}
		}
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()
	return postsS, nil
}

func (r *repoPostgres) UpdateThreadInfo(upThread models.Thread) (models.Thread, models.StatusCode) {
	threadS := models.Thread{}
	if upThread.Slug == "" {
		rowQuery := fmt.Sprintf(UpdateThread, `id=$3`)
		row := r.Conn.QueryRow(context.Background(), rowQuery, upThread.Title, upThread.Message, upThread.ID)
		err := row.Scan(&threadS.ID, &threadS.Title, &threadS.Author, &threadS.Created,
			&threadS.Forum, &threadS.Message, &threadS.Slug, &threadS.Votes)
		if err != nil {
			return models.Thread{}, models.NotFound
		}
	} else {
		rowQuery := fmt.Sprintf(UpdateThread, `slug=$3`)
		row := r.Conn.QueryRow(context.Background(), rowQuery, upThread.Title, upThread.Message, upThread.Slug)
		err := row.Scan(&threadS.ID, &threadS.Title, &threadS.Author, &threadS.Created,
			&threadS.Forum, &threadS.Message, &threadS.Slug, &threadS.Votes)
		if err != nil {
			return models.Thread{}, models.NotFound
		}
	}
	return threadS, models.Okey
}

func (r *repoPostgres) GetPostsFlat(limit string, since string, desc string, ID int) ([]models.Post, models.StatusCode) {
	manyPosts := make([]models.Post, 0)
	if since != "" {
		if desc != "" {
			rows, err := r.Conn.Query(context.Background(), SelectPostSinceDescNotNil, ID, limit)
			if err != nil {
				return manyPosts, models.InternalError
			}
			defer rows.Close()
			for rows.Next() {
				onePost := models.Post{}
				err := rows.Scan(&onePost.ID, &onePost.Author, &onePost.Created, &onePost.Forum, &onePost.IsEdited, &onePost.Message, &onePost.Parent, &onePost.Thread)
				if err != nil {
					return manyPosts, models.InternalError
				}
				manyPosts = append(manyPosts, onePost)
			}
		} else {
			rows, err := r.Conn.Query(context.Background(), SelectPostSinceDescNil, ID, limit)
			if err != nil {
				return manyPosts, models.InternalError
			}
			defer rows.Close()
			for rows.Next() {
				onePost := models.Post{}
				err := rows.Scan(&onePost.ID, &onePost.Author, &onePost.Created, &onePost.Forum, &onePost.IsEdited, &onePost.Message, &onePost.Parent, &onePost.Thread)
				if err != nil {
					return manyPosts, models.InternalError
				}
				manyPosts = append(manyPosts, onePost)
			}
		}
	} else {
		if desc != "" {
			rows, err := r.Conn.Query(context.Background(), SelectPostDescNotNil, ID, since, limit)
			if err != nil {
				return manyPosts, models.InternalError
			}
			defer rows.Close()
			for rows.Next() {
				onePost := models.Post{}
				err := rows.Scan(&onePost.ID, &onePost.Author, &onePost.Created, &onePost.Forum, &onePost.IsEdited, &onePost.Message, &onePost.Parent, &onePost.Thread)
				if err != nil {
					return manyPosts, models.InternalError
				}
				manyPosts = append(manyPosts, onePost)
			}
		} else {
			rows, err := r.Conn.Query(context.Background(), SelectPostDescNil, ID, since, limit)
			if err != nil {
				return manyPosts, models.InternalError
			}
			defer rows.Close()
			for rows.Next() {
				onePost := models.Post{}
				err := rows.Scan(&onePost.ID, &onePost.Author, &onePost.Created, &onePost.Forum, &onePost.IsEdited, &onePost.Message, &onePost.Parent, &onePost.Thread)
				if err != nil {
					return manyPosts, models.InternalError
				}
				manyPosts = append(manyPosts, onePost)
			}
		}
	}
	return manyPosts, models.Okey
}

func (r *repoPostgres) GetPostsTree(limit string, since string, desc string, ID int) ([]models.Post, models.StatusCode) {
	manyPosts := make([]models.Post, 0)
	if since == "" {
		if desc != "" {
			rows, err := r.Conn.Query(context.Background(), SelectPostTreeSinceDescNotNil, ID, limit)
			if err != nil {
				return manyPosts, models.InternalError
			}
			defer rows.Close()
			for rows.Next() {
				onePost := models.Post{}
				err = rows.Scan(&onePost.ID, &onePost.Author, &onePost.Created, &onePost.Forum, &onePost.IsEdited, &onePost.Message, &onePost.Parent, &onePost.Thread)
				if err != nil {
					return manyPosts, models.InternalError
				}
				manyPosts = append(manyPosts, onePost)
			}
		} else {
			rows, err := r.Conn.Query(context.Background(), SelectPostTreeSinceDescNil, ID, limit)
			if err != nil {
				return manyPosts, models.InternalError
			}
			defer rows.Close()
			for rows.Next() {
				onePost := models.Post{}
				err = rows.Scan(&onePost.ID, &onePost.Author, &onePost.Created, &onePost.Forum, &onePost.IsEdited, &onePost.Message, &onePost.Parent, &onePost.Thread)
				if err != nil {
					return manyPosts, models.InternalError
				}
				manyPosts = append(manyPosts, onePost)
			}
		}
	} else {
		if desc == "" {
			rows, err := r.Conn.Query(context.Background(), SelectPostTreeDescNotNil, ID, since, limit)
			if err != nil {
				return manyPosts, models.InternalError
			}
			defer rows.Close()
			for rows.Next() {
				onePost := models.Post{}
				err = rows.Scan(&onePost.ID, &onePost.Author, &onePost.Created, &onePost.Forum, &onePost.IsEdited, &onePost.Message, &onePost.Parent, &onePost.Thread)
				if err != nil {
					return manyPosts, models.InternalError
				}
				manyPosts = append(manyPosts, onePost)
			}
		} else {
			rows, err := r.Conn.Query(context.Background(), SelectPostTreeDescNil, ID, since, limit)
			if err != nil {
				return manyPosts, models.InternalError
			}
			defer rows.Close()
			for rows.Next() {
				onePost := models.Post{}
				err = rows.Scan(&onePost.ID, &onePost.Author, &onePost.Created, &onePost.Forum, &onePost.IsEdited, &onePost.Message, &onePost.Parent, &onePost.Thread)
				if err != nil {
					return manyPosts, models.InternalError
				}
				manyPosts = append(manyPosts, onePost)
			}
		}
	}
	return manyPosts, models.Okey
}

func (r *repoPostgres) GetPostsParent(limit string, since string, desc string, ID int) ([]models.Post, models.StatusCode) {
	manyPosts := make([]models.Post, 0)

	if since == "" {
		if desc != "" {
			rows, err := r.Conn.Query(context.Background(), SelectPostParentSinceDescNotNil, ID, limit)
			if err != nil {
				return manyPosts, models.InternalError
			}
			defer rows.Close()
			for rows.Next() {
				onePost := models.Post{}
				err = rows.Scan(&onePost.ID, &onePost.Author, &onePost.Created, &onePost.Forum, &onePost.IsEdited, &onePost.Message, &onePost.Parent, &onePost.Thread)
				if err != nil {
					return manyPosts, models.InternalError
				}
				manyPosts = append(manyPosts, onePost)
			}
		} else {
			rows, err := r.Conn.Query(context.Background(), SelectPostParentSinceDescNil, ID, limit)
			if err != nil {
				return manyPosts, models.InternalError
			}
			defer rows.Close()
			for rows.Next() {
				onePost := models.Post{}
				err = rows.Scan(&onePost.ID, &onePost.Author, &onePost.Created, &onePost.Forum, &onePost.IsEdited, &onePost.Message, &onePost.Parent, &onePost.Thread)
				if err != nil {
					return manyPosts, models.InternalError
				}
				manyPosts = append(manyPosts, onePost)
			}
		}
	} else {
		if desc != "" {
			rows, err := r.Conn.Query(context.Background(), SelectPostParentDescNotNil, ID, since, limit)
			if err != nil {
				return manyPosts, models.InternalError
			}
			defer rows.Close()
			for rows.Next() {
				onePost := models.Post{}
				err = rows.Scan(&onePost.ID, &onePost.Author, &onePost.Created, &onePost.Forum, &onePost.IsEdited, &onePost.Message, &onePost.Parent, &onePost.Thread)
				if err != nil {
					return manyPosts, models.InternalError
				}
				manyPosts = append(manyPosts, onePost)
			}
		} else {
			rows, err := r.Conn.Query(context.Background(), SelectPostParentDescNil, ID, since, limit)
			if err != nil {
				return manyPosts, models.InternalError
			}
			defer rows.Close()
			for rows.Next() {
				onePost := models.Post{}
				err = rows.Scan(&onePost.ID, &onePost.Author, &onePost.Created, &onePost.Forum, &onePost.IsEdited, &onePost.Message, &onePost.Parent, &onePost.Thread)
				if err != nil {
					return manyPosts, models.InternalError
				}
				manyPosts = append(manyPosts, onePost)
			}
		}
	}
	return manyPosts, models.Okey
}

func (r *repoPostgres) InVoted(vote models.Vote) error {
	_, err := r.Conn.Exec(context.Background(), InsertVote, vote.Nickname, vote.Voice, vote.Thread)
	if err != nil {
		return err
	}
	return nil
}

func (r *repoPostgres) UpVote(vote models.Vote) (models.Vote, error) {
	_, err := r.Conn.Exec(context.Background(), UpdateVote, vote.Nickname, vote.Voice, vote.Thread)
	if err != nil {
		return models.Vote{}, err
	}
	return vote, nil
}
