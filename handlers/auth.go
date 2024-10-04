package handlers

import (
	"log"
	"net/http"

	"github.com/Open-Code-Zone/cms/views/pages"
	"github.com/markbates/goth/gothic"
)

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	pages.LoginPage().Render(r.Context(), w)
}

func (h *Handler) HandleProviderLogin(w http.ResponseWriter, r *http.Request) {
	u, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		gothic.BeginAuthHandler(w, r)
	} else {
    log.Println(u)
		http.Redirect(w, r, "/blog-post", http.StatusSeeOther)
	}
}

func (h *Handler) HandleAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	u, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
    log.Println(err)
    return
	}

  err = h.auth.StoreUserSession(w, r, u)

  w.Header().Set("Location", "/blog-post")
  w.WriteHeader(http.StatusTemporaryRedirect)
}
