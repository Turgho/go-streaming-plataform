package handler

import "net/http"

type VideoRouter struct {
	mux *http.ServeMux
}

func NewVideoRoute(mux *http.ServeMux) *VideoRouter {
	return &VideoRouter{mux: mux}
}

func (r *VideoRouter) RegisterRoutes(handler *VideoHandler) {
	r.mux.Handle("/videos/", WithHeaders(handler.ServeFiles()))
}
