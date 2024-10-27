package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Open-Code-Zone/cms/services/auth"
	"github.com/Open-Code-Zone/cms/views/pages"
	"github.com/markbates/goth/gothic"
)

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	_, err := auth.GetSessionUser(r)
	if err != nil {
		log.Println(err)
		pages.LoginPage().Render(r.Context(), w)
	}

	w.Header().Set("Location", "/blog-post")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) ProviderLogin(w http.ResponseWriter, r *http.Request) {
	if u, err := gothic.CompleteUserAuth(w, r); err == nil {
		log.Printf("User already authenticated! %v", u)

		pages.LoginPage().Render(r.Context(), w)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (h *Handler) AuthCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	err = auth.StoreUserSession(w, r, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/blog-post")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("Logging out...")

	err := gothic.Logout(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	auth.RemoveUserSession(w, r)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
