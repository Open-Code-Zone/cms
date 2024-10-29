package handlers

import (
	"net/http"

	"github.com/Open-Code-Zone/cms/store"
	"github.com/Open-Code-Zone/cms/utils"
)

type Handler struct {
	store        *store.Storage
	githubClient *utils.GitHubClient
}

func New(store *store.Storage, githubClient *utils.GitHubClient) *Handler {
	return &Handler{
		store:        store,
		githubClient: githubClient,
	}
}

func (h *Handler) PingIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All good server is up and running!"))
	w.WriteHeader(http.StatusOK)
	return
}
