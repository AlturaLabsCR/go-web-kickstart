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

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var registerRequest struct {
		Email string `json:"email"`
	}

	data, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	json.Unmarshal(data, &registerRequest)

	userID, err := database.InsertUser(h.DB(), r.Context(), registerRequest.Email)
	if err != nil {
		h.Log().Error("error registering user", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		if errors.Is(err, database.ErrDuplicateEmail) {
			fmt.Fprintln(w, "user already exists")
		} else {
			fmt.Fprintln(w, "error registering user")
		}
		return
	}

	fmt.Fprintln(w, "registered user", userID)
}

func (h *Handler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	content := templates.Register(h.Translator(r))
	templates.Base(content).Render(r.Context(), w)
}
