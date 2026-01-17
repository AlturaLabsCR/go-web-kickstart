package handler

import (
	"net/http"

	"app/config/routes"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.Sess().Revoke(w, r); err != nil {
		h.Log().Error("error revoking sessions", "error", err)
		return
	}
	http.Redirect(w, r, routes.Map[routes.Login], http.StatusSeeOther)
}
