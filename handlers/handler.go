package handlers

import "githum.com/Open-Code-Zone/cms/store"

type Handler struct {
	store *store.Storage
}

func New(store *store.Storage) *Handler {
	return &Handler{
		store: 		store,
	}
}
