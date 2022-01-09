package user

import (
	"fmt"
	"net/http"

	"github.com/iandjx/go-oauth-2/pkg/httputil"
)

type Handler struct {
	user Service
}

type CreateParam struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewHandler(us Service) *Handler {
	return &Handler{us}
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
