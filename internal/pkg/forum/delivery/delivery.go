package delivery

import (
	"forumI/internal/models"
	"forumI/internal/pkg/forum"
	"forumI/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"net/http"
)

type Handler struct {
	uc forum.UseCase
}

func NewForumHandler(ForumUseCase forum.UseCase) *Handler {
	return &Handler{uc: ForumUseCase}
}

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
