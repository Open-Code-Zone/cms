package handlers

import (
	"github.com/Open-Code-Zone/cms/services/auth"
	"github.com/Open-Code-Zone/cms/store"
)

type Handler struct {
	store *store.Storage
	auth  *auth.AuthService
}

func New(store *store.Storage, auth *auth.AuthService) *Handler {
	return &Handler{
		store: store,
		auth:  auth,
	}
}
