package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"app/database"
	"app/templates"
)

func (h *Handler) RegisterOwner(w http.ResponseWriter, r *http.Request) {
	var registerRequest struct {
		Email string `json:"email"`
	}

	data, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	json.Unmarshal(data, &registerRequest)

	ownerID, err := database.InsertOwner(h.DB(), r.Context(), registerRequest.Email)
	if err != nil {
		h.Log().Error("error registering owner", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		if errors.Is(err, database.ErrDuplicateEmail) {
			fmt.Fprintln(w, "owner already exists")
		} else {
			fmt.Fprintln(w, "error registering owner")
		}
		return
	}

	fmt.Fprintln(w, "registered owner", ownerID)
}

func (h *Handler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	content := templates.Register(h.Translator(r))
	templates.Base(content).Render(r.Context(), w)
}
