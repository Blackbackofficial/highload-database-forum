package forum

import "forumI/internal/models"

type UseCase interface {
	Forum(forum models.Forum) (models.Forum, models.StatusCode)
	GetForum(forum models.Forum) (models.Forum, models.StatusCode)
}

type Repository interface {
	GetUser(forum models.Forum) (models.User, models.StatusCode)
	InForum(forum models.Forum) error
	GetForum(forum models.Forum) (models.Forum, models.StatusCode)
}
