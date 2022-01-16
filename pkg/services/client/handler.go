package client

import (
	"net/http"

	"github.com/iandjx/go-oauth-2/pkg/httputil"
)

type Handler struct {
	client Service
}

type CreateParam struct {
	Name         string   `json:"name"`
	RedirectURLs []string `json:"redirect_urls"`
	Scope        []string `json:"scopes"`
}

func NewHandler(cs Service) *Handler {
	return &Handler{cs}
}

func (h *Handler) CreateClient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var p CreateParam
		if err := httputil.DecodeJSON(r, &p); err != nil {
			httputil.Error400(w, err)
			return
		}

		out, err := h.client.CreateClient(ctx, p)
		if err != nil {
			httputil.Error400(w, err)
			return
		}
		if err := httputil.EncodeJSON(w, out, http.StatusOK); err != nil {
			httputil.Error500(w, err)
		}

	}
}

// TODO add handler to get client
// TODO add handler to delete client
