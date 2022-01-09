package usecase

import (
	"forumI/internal/models"
	"forumI/internal/pkg/forum"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"strconv"
)

type UseCase struct {
	repo forum.Repository
}

func NewRepoUsecase(repo forum.Repository) forum.UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) Forum(forum models.Forum) (models.Forum, models.StatusCode) {
	user, status := u.repo.GetUser(forum.User)
	if status != models.Okey {
		return models.Forum{}, status
	}
	forum.User = user.NickName

	errI := u.repo.InForum(forum)
	if errI != nil {
		if pgError, ok := errI.(pgx.PgError); ok && pgError.Code == "23503" {
			return models.Forum{}, models.NotFound
		}
		if pgError, ok := errI.(pgx.PgError); ok && pgError.Code == "23505" {
			forumM, _ := u.repo.GetForum(forum.Slug)
			return forumM, models.Conflict
		}
		return models.Forum{}, models.InternalError
	}

	forum.Posts = 0
	forum.Threads = 0
	return forum, models.Created
}

func (u *UseCase) GetForum(forum models.Forum) (models.Forum, models.StatusCode) {
	return u.repo.GetForum(forum.Slug)
}

func (u *UseCase) CreateThreadsForum(thread models.Thread) (models.Thread, models.StatusCode) {
	forumS, status := u.repo.GetForum(thread.Forum)
	if status != models.Okey {
		return models.Thread{}, status
	}

	user, status := u.repo.GetUser(thread.Author)
	if status != models.Okey {
		return models.Thread{}, status
	}
	slug := uuid.New().String()
	thread.Slug = slug
	thread.Author = user.NickName
	thread.Forum = forumS.Slug

	thread, errI := u.repo.InThread(thread)
	if errI != nil {
		if pgError, ok := errI.(pgx.PgError); ok && pgError.Code == "23505" {
			threadM, _ := u.repo.GetThreadSlug(slug)
			return threadM, models.Conflict
		}
		return models.Thread{}, models.InternalError
	}
	return thread, models.Created
}

func (u *UseCase) GetUsersOfForum(forum models.Forum, limit string, since string, desc string) ([]models.User, models.StatusCode) {
	_, status := u.repo.GetForum(forum.Slug)
	if status != models.Okey {
		return nil, status
	}

	users, status := u.repo.GetUsersOfForum(forum, limit, since, desc)
	if status != models.Okey {
		return nil, status
	}
	return users, models.Okey
}

func (u *UseCase) GetThreadsOfForum(forum models.Forum, limit string, since string, desc string) ([]models.Thread, models.StatusCode) {
	_, status := u.repo.GetForum(forum.Slug)
	if status != models.Okey {
		return nil, status
	}

	threads, status := u.repo.GetThreadsOfForum(forum, limit, since, desc)
	if status != models.Okey {
		return nil, status
	}
	return threads, models.Okey
}

func (u *UseCase) GetFullPostInfo(posts models.PostFull, related []string) (models.PostFull, models.StatusCode) {
	return u.repo.GetFullPostInfo(posts, related)
}

func (u *UseCase) UpdatePostInfo(postUpdate models.PostUpdate) (models.Post, models.StatusCode) {
	postOne := models.Post{ID: postUpdate.ID}
	postOne, status := u.repo.UpdatePostInfo(postOne, postUpdate)
	if status != models.Okey {
		return models.Post{}, models.NotFound
	}
	return postOne, models.Okey
}

func (u *UseCase) GetClear() models.StatusCode {
	return u.repo.GetClear()
}

func (u *UseCase) GetStatus() models.Status {
	return u.repo.GetStatus()
}

func (u *UseCase) CheckThreadIdOrSlug(slugOrId string) (models.Thread, models.StatusCode) {
	isInt, err := strconv.Atoi(slugOrId)
	if err != nil {
		threadS, status := u.repo.GetThreadSlug(slugOrId)
		if status != models.Okey {
			return models.Thread{}, status
		}
		return threadS, models.Okey
	} else {
		threadS, status := u.repo.GetIdThread(isInt)
		if status != models.Okey {
			return models.Thread{}, status
		}
		return threadS, models.Okey
	}
}

func (u *UseCase) CreatePosts(inPosts []models.Post, thread models.Thread) ([]models.Post, models.StatusCode) {
	posts := make([]models.Post, 0)
	posts, err := u.repo.InPosts(inPosts, thread)
	if err != nil {
		if pgError, ok := err.(pgx.PgError); ok && pgError.Code == "23503" {
			return nil, models.NotFound
		} else {
			return nil, models.Conflict
		}
	}
	return posts, models.Created
}

func (u *UseCase) UpdateThreadInfo(slugOrId string, upThread models.Thread) (models.Thread, models.StatusCode) {
	isInt, err := strconv.Atoi(slugOrId)
	if err != nil {
		upThread.Slug = slugOrId
	} else {
		upThread.ID = isInt
	}
	return u.repo.UpdateThreadInfo(upThread)
}

func (u *UseCase) GetPostOfThread(limit string, since string, desc string, sort string, ID int) ([]models.Post, models.StatusCode) {
	switch sort {
	case "flat":
		return u.repo.GetPostsFlat(limit, since, desc, ID)
	case "tree":
		return u.repo.GetPostsTree(limit, since, desc, ID)
	case "parent_tree":
		return u.repo.GetPostsParent(limit, since, desc, ID)
	default:
		return u.repo.GetPostsFlat(limit, since, desc, ID)
	}
}
