package delivery

import (
	"forumI/internal/models"
	"forumI/internal/pkg/forum"
	"forumI/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	uc forum.UseCase
}

func NewForumHandler(ForumUseCase forum.UseCase) *Handler {
	return &Handler{uc: ForumUseCase}
}

// CreateForum /forum/create
func (h Handler) CreateForum(w http.ResponseWriter, r *http.Request) {
	forumS := models.Forum{}
	err := easyjson.UnmarshalFromReader(r.Body, &forumS)
	if err != nil {
		utils.Response(w, models.InternalError, nil)
		return
	}

	forumS, status := h.uc.Forum(forumS)
	utils.Response(w, status, forumS)
}

// ForumInfo /forum/{slug}/details
func (h Handler) ForumInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, found := vars["slug"]
	if !found {
		utils.Response(w, models.NotFound, nil)
		return
	}
	forumS := models.Forum{Slug: slug}
	forumS, status := h.uc.GetForum(forumS)
	utils.Response(w, status, forumS)
}

// CreateThreadsForum /forum/{slug}/create
func (h Handler) CreateThreadsForum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, found := vars["slug"]
	if !found {
		utils.Response(w, models.NotFound, nil)
		return
	}

	threadS := models.Thread{}
	err := easyjson.UnmarshalFromReader(r.Body, &threadS)
	if err != nil {
		utils.Response(w, models.InternalError, nil)
		return
	}
	threadS.Forum = slug

	threadS, status := h.uc.CreateThreadsForum(threadS)
	utils.Response(w, status, threadS)
}

// GetUsersForum /forum/{slug}/users
func (h Handler) GetUsersForum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, found := vars["slug"]
	if !found {
		utils.Response(w, models.NotFound, nil)
		return
	}
	limit := ""
	since := ""
	desc := ""

	query := r.URL.Query()
	if limits := query["limit"]; len(limits) > 0 {
		limit = limits[0]
	}
	if sinces := query["since"]; len(sinces) > 0 {
		since = sinces[0]
	}
	if descs := query["desc"]; len(descs) > 0 {
		desc = descs[0]
	}

	forumS := models.Forum{Slug: slug}

	users, status := h.uc.GetUsersOfForum(forumS, limit, since, desc)
	utils.Response(w, status, users)
}

// GetThreadsForum /forum/{slug}/threads
func (h Handler) GetThreadsForum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, found := vars["slug"]
	if !found {
		utils.Response(w, models.NotFound, nil)
		return
	}
	limit := ""
	since := ""
	desc := ""

	query := r.URL.Query()
	if limits := query["limit"]; len(limits) > 0 {
		limit = limits[0]
	}
	if sinces := query["since"]; len(sinces) > 0 {
		since = sinces[0]
	}
	if descs := query["desc"]; len(descs) > 0 {
		desc = descs[0]
	}
	forumS := models.Forum{Slug: slug}

	users, status := h.uc.GetThreadsOfForum(forumS, limit, since, desc)
	utils.Response(w, status, users)
}

// GetPostInfo /post/{id}/details
func (h Handler) GetPostInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idV := vars["id"]
	idV, found := vars["id"]
	if !found {
		utils.Response(w, models.NotFound, nil)
		return
	}

	id, _ := strconv.Atoi(idV)
	query := r.URL.Query()

	var related []string
	if relateds := query["related"]; len(relateds) > 0 {
		related = strings.Split(relateds[0], ",")
	}

	postFull := models.PostFull{}
	err := easyjson.UnmarshalFromReader(r.Body, &postFull)
	if err != nil {
		utils.Response(w, models.InternalError, nil)
		return
	}

	postFull.Post.ID = id
	finalPostF, status := h.uc.GetFullPostInfo(postFull, related)
	utils.Response(w, status, finalPostF)
}
