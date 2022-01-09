package forum

import "forumI/internal/models"

type UseCase interface {
	Forum(forum models.Forum) (models.Forum, models.StatusCode)
	GetForum(forum models.Forum) (models.Forum, models.StatusCode)
	CreateThreadsForum(thread models.Thread) (models.Thread, models.StatusCode)
	GetUsersOfForum(forum models.Forum, limit string, since string, desc string) ([]models.User, models.StatusCode)
	GetThreadsOfForum(forum models.Forum, limit string, since string, desc string) ([]models.Thread, models.StatusCode)
	GetFullPostInfo(posts models.PostFull, related []string) (models.PostFull, models.StatusCode)
}

type Repository interface {
	GetUser(name string) (models.User, models.StatusCode)
	InForum(forum models.Forum) error
	GetForum(slug string) (models.Forum, models.StatusCode)
	InThread(thread models.Thread) (models.Thread, error)
	GetThreadSlug(slug string) (models.Thread, models.StatusCode)
	GetUsersOfForum(forum models.Forum, limit string, since string, desc string) ([]models.User, models.StatusCode)
	GetThreadsOfForum(forum models.Forum, limit string, since string, desc string) ([]models.Thread, models.StatusCode)
	GetFullPostInfo(posts models.PostFull, related []string) (models.PostFull, models.StatusCode)
	GetIdThread(id int) (models.Thread, models.StatusCode)
}
