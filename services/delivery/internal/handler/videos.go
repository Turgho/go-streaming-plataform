package handler

import "net/http"

type VideoHandler struct {
	dir string
}

func NewVideoHandler(dir string) *VideoHandler {
	return &VideoHandler{dir: dir}
}

func (h *VideoHandler) ServeFiles() http.Handler {
	fs := http.FileServer(http.Dir(h.dir))
	return http.StripPrefix("/videos/", fs)
}
