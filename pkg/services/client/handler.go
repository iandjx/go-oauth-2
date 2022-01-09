package client

import (
	"fmt"
	"net/http"

	"github.com/iandjx/go-oauth-2/pkg/httputil"
)

type Handler struct {
	client Service
}

type CreateParam struct {
	Name         string   `json:"name"`
	RedirectURLs []string `json:"redirect_urls"`
	Email        string   `json:"email"`
	Password     string   `json:"password"`
}

func NewHandler(cs Service) *Handler {
	return &Handler{cs}
}

func (h *Handler) RegisterUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hit")

		var p CreateParam
		if err := httputil.DecodeJSON(r, &p); err != nil {
			httputil.Error400(w, err)
			return
		}

		out, err := h.user.Register(p)
		if err != nil {
			httputil.Error400(w, err)
			return
		}
		authCookie := &http.Cookie{Name: "auth", Value: out, HttpOnly: true}

		http.SetCookie(w, authCookie)

		w.Write([]byte("OK"))
	}
}
