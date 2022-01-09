package delivery

import (
	"encoding/json"
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
func (h *Handler) CreateForum(w http.ResponseWriter, r *http.Request) {
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
func (h *Handler) ForumInfo(w http.ResponseWriter, r *http.Request) {
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
func (h *Handler) CreateThreadsForum(w http.ResponseWriter, r *http.Request) {
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
func (h *Handler) GetUsersForum(w http.ResponseWriter, r *http.Request) {
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
func (h *Handler) GetThreadsForum(w http.ResponseWriter, r *http.Request) {
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
func (h *Handler) GetPostInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
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

// UpdatePostInfo /post/{id}/details
func (h *Handler) UpdatePostInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids, found := vars["id"]
	if !found {
		utils.Response(w, models.NotFound, nil)
		return
	}

	postUpdate := models.PostUpdate{}
	err := easyjson.UnmarshalFromReader(r.Body, &postUpdate)
	if err != nil {
		utils.Response(w, models.InternalError, nil)
		return
	}
	id, err := strconv.Atoi(ids)

	if err == nil {
		postUpdate.ID = id
	}

	finalPostU, status := h.uc.UpdatePostInfo(postUpdate)
	utils.Response(w, status, finalPostU)
}

// GetClear /service/clear
func (h *Handler) GetClear(w http.ResponseWriter, _ *http.Request) {
	status := h.uc.GetClear()
	utils.Response(w, status, nil)
}

// GetStatus /service/status
func (h *Handler) GetStatus(w http.ResponseWriter, _ *http.Request) {
	statusS := h.uc.GetStatus()
	utils.Response(w, models.Okey, statusS)
}

// CreatePosts /thread/{slug_or_id}/create
func (h *Handler) CreatePosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId, found := vars["slug_or_id"]
	if !found {
		utils.Response(w, models.NotFound, nil)
		return
	}

	var posts []models.Post
	thread, status := h.uc.CheckThreadIdOrSlug(slugOrId)
	if status != models.Okey {
		utils.Response(w, status, nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&posts)
	if err != nil {
		utils.Response(w, models.InternalError, nil)
		return
	}

	if len(posts) == 0 {
		utils.Response(w, models.Created, []byte("[]"))
		return
	}

	createPosts, status := h.uc.CreatePosts(posts, thread)
	utils.Response(w, status, createPosts)
}

func (h Handler) GetThreadInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId, found := vars["slug_or_id"]
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
	finalThread, status := h.uc.CheckThreadIdOrSlug(slugOrId)
	utils.Response(w, status, finalThread)
}
