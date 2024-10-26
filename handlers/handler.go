package handlers

import (
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
