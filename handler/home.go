package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "hello, world\n")
}
