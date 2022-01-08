package usecase

import (
	"forumI/internal/models"
	"forumI/internal/pkg/forum"
	"github.com/jackc/pgx"
)

type UseCase struct {
	repo forum.Repository
}

func NewRepoUsecase(repo forum.Repository) forum.UseCase {
	return &UseCase{repo: repo}
}

func (u UseCase) Forum(forum models.Forum) (models.Forum, models.StatusCode) {
	user, err := u.repo.GetUser(forum)
	if err != models.Okey {
		return models.Forum{}, err
	}
	forum.User = user.NickName

	errI := u.repo.InForum(forum)
	if errI != nil {
		if pgError, ok := errI.(pgx.PgError); ok && pgError.Code == "23503" {
			return models.Forum{}, models.NotFound
		}
		if pgError, ok := errI.(pgx.PgError); ok && pgError.Code == "23505" {
			forumM, _ := u.repo.GetForum(forum)
			return forumM, models.Conflict
		}
		return models.Forum{}, models.InternalError
	}

	forum.Posts = 0
	forum.Threads = 0

	return forum, models.Created
}

func (u UseCase) GetForum(forum models.Forum) (models.Forum, models.StatusCode) {
	return u.repo.GetForum(forum)
}
