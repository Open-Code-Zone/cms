package handlers

import (
	"net/http"

	"githum.com/Open-Code-Zone/cms/views"
)

func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) {
	views.Home().Render(r.Context(), w)
}
