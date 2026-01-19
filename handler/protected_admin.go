package handler

import (
	"net/http"

	"app/config/routes"
	"app/database/models"
	"app/templates/base"
	"app/templates/protected"
)

func (h *Handler) ProtectedAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionData, ok := h.Sess().Data(ctx)
	if !ok {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	sessionAttrs, ok := h.Sess().Attrs(ctx)
	if !ok {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	tr := h.Tr(r)

	userMeta, err := h.DB().Querier().GetUserMeta(ctx, sessionData.UserID)
	if err != nil {
		h.Log().Error("error getting user meta", "error", err)
		http.Error(w, "error getting user meta", http.StatusInternalServerError)
		return
	}

	if !models.HasPermission(userMeta.Perms, "perm.change_name") {
		h.Log().Error("doesnt have permission to change name")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	allUsers, err := h.DB().Querier().GetUsersMeta(ctx)
	if err != nil {
		h.Log().Error("error getting all users", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	main := protected.ProtectedMain(tr, protected.ProtectedParams{
		User:     userMeta,
		Attrs:    sessionAttrs,
		Data:     sessionData,
		AllUsers: allUsers,
		Active:   r.URL.Path,
	})

	head := base.HeadParams{
		Subtitle:    tr("nav.account"),
		LoadJS:      true,
		RobotsIndex: false,
	}

	asideParams := protected.AsideParams{
		Active:  r.URL.Path,
		IsAdmin: models.HasPermission(userMeta.Perms, "perm.admin"),
	}

	aside := protected.Aside(tr, asideParams)

	page := base.PageParams{
		Head: head,
		Body: base.BodyParams{
			Content: main,
			Aside:   aside,
			Active:  routes.Map[routes.Protected],
		},
	}

	if err := base.Page(tr, page).Render(ctx, w); err != nil {
		h.Log().Error("error rendering template", "error", err)
	}
}

func (h *Handler) ProtectedDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionData, ok := h.Sess().Data(ctx)
	if !ok {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	tr := h.Tr(r)

	userMeta, err := h.DB().Querier().GetUserMeta(ctx, sessionData.UserID)
	if err != nil {
		h.Log().Error("error getting user meta", "error", err)
		http.Error(w, "error getting user meta", http.StatusInternalServerError)
		return
	}

	if !models.HasPermission(userMeta.Perms, "perm.manage_users") {
		h.Log().Error("doesnt have permission to manage users")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	idToDelete := r.PathValue("id")
	if idToDelete == userMeta.ID {
		h.Log().Error("error deleting user", "error", "cannot delete oneself")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.DB().Querier().DelUser(ctx, idToDelete); err != nil {
		h.Log().Error("error deleting user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	allUsers, err := h.DB().Querier().GetUsersMeta(ctx)
	if err != nil {
		h.Log().Error("error getting all users", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := protected.ProtectedAllUsers(tr, allUsers, userMeta.ID).Render(ctx, w); err != nil {
		h.Log().Error("error rendering template", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
