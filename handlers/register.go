package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (h *Handler) RegisterOwner(w http.ResponseWriter, r *http.Request) {
	var registerRequest struct {
		Email string `json:"email"`
	}

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	json.Unmarshal(body, &registerRequest)

	h.DB().InsertOwner(r.Context(), registerRequest.Email)

	fmt.Fprintln(w, "registered owner", registerRequest.Email)
}
