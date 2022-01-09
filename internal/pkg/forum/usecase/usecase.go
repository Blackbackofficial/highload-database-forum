package usecase

import (
	"forumI/internal/models"
	"forumI/internal/pkg/forum"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

type UseCase struct {
	repo forum.Repository
}

func NewRepoUsecase(repo forum.Repository) forum.UseCase {
	return &UseCase{repo: repo}
}

func (u UseCase) Forum(forum models.Forum) (models.Forum, models.StatusCode) {
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

func (u UseCase) GetForum(forum models.Forum) (models.Forum, models.StatusCode) {
	return u.repo.GetForum(forum.Slug)
}

func (u UseCase) CreateThreadsForum(thread models.Thread) (models.Thread, models.StatusCode) {
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

func (u UseCase) GetUsersOfForum(forum models.Forum, limit string, since string, desc string) ([]models.User, models.StatusCode) {
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

func (u UseCase) GetThreadsOfForum(forum models.Forum, limit string, since string, desc string) ([]models.Thread, models.StatusCode) {
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
