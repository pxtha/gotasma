package api

import (
	"io"
	"net/http"
)

type (
	IndexHandler struct {
	}
)

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf8")
	io.WriteString(w, `
		<h5 style="color: #111; font-family: 'Helvetica Neue', sans-serif; font-size: 100px; font-weight: bold; letter-spacing: -1px; line-height: 1; text-align: center; ">Gotasma</h2>
		<h5 style="color: #111; font-family: 'Open Sans', sans-serif; font-size: 30px; font-weight: 300; line-height: 32px; margin: 0 0 72px; text-align: center; ">Welcome to my Project</h3>
	`)
}
