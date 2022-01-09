package forum

import "forumI/internal/models"

type UseCase interface {
	Forum(forum models.Forum) (models.Forum, models.StatusCode)
	GetForum(forum models.Forum) (models.Forum, models.StatusCode)
	CreateThreadsForum(thread models.Thread) (models.Thread, models.StatusCode)
	GetUsersOfForum(forum models.Forum, limit string, since string, desc string) ([]models.User, models.StatusCode)
	GetThreadsOfForum(forum models.Forum, limit string, since string, desc string) ([]models.Thread, models.StatusCode)
	GetFullPostInfo(posts models.PostFull, related []string) (models.PostFull, models.StatusCode)
	UpdatePostInfo(postUpdate models.PostUpdate) (models.Post, models.StatusCode)
	GetClear() models.StatusCode
	GetStatus() models.Status
	CheckThreadIdOrSlug(slugOrId string) (models.Thread, models.StatusCode)
	CreatePosts(createPosts []models.Post, thread models.Thread) ([]models.Post, models.StatusCode)
	UpdateThreadInfo(slugOrId string, upThread models.Thread) (models.Thread, models.StatusCode)
	GetPostOfThread(limit string, since string, desc string, sort string, ID int) ([]models.Post, models.StatusCode)
	Voted(vote models.Vote, thread models.Thread) (models.Thread, models.StatusCode)
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
	UpdatePostInfo(post models.Post, postUpdate models.PostUpdate) (models.Post, models.StatusCode)
	GetClear() models.StatusCode
	GetStatus() models.Status
	InPosts(posts []models.Post, thread models.Thread) ([]models.Post, error)
	UpdateThreadInfo(upThread models.Thread) (models.Thread, models.StatusCode)
	GetPostsFlat(limit string, since string, desc string, ID int) ([]models.Post, models.StatusCode)
	GetPostsTree(limit string, since string, desc string, ID int) ([]models.Post, models.StatusCode)
	GetPostsParent(limit string, since string, desc string, ID int) ([]models.Post, models.StatusCode)
	InVoted(vote models.Vote) error
	UpVote(vote models.Vote) (models.Vote, error)
}
