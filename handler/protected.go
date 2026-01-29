package handler

import (
	"net/http"

	"app/config/routes"
	"app/database/models"
	"app/i18n"
	"app/templates/base"
	"app/templates/protected"
)

func (h *Handler) ProtectedPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionData, ok := h.Sess().Data(ctx)
	if !ok {
		h.Log().Error("session not found")
		http.Redirect(w, r, routes.Map[routes.Login], http.StatusSeeOther)
		return
	}

	sessionAttrs, ok := h.Sess().Attrs(ctx)
	if !ok {
		h.Log().Error("session not found")
		http.Redirect(w, r, routes.Map[routes.Login], http.StatusSeeOther)
		return
	}

	tr := h.Tr(r)

	userMeta, err := h.DB().Querier().GetUserMeta(ctx, sessionData.UserID)
	if err != nil {
		h.Log().Error("error getting user meta", "error", err)
		http.Redirect(w, r, routes.Map[routes.Login], http.StatusSeeOther)
		return
	}

	sessions, err := h.Sess().AttrsByUser(ctx, sessionData.UserID)
	if err != nil {
		h.Log().Error("error getting sessions attrs", "error", err)
		http.Error(w, "error getting session attrs", http.StatusInternalServerError)
		return
	}

	locale := "es"

	langs := i18n.RequestLanguages(r)
	if len(langs) > 0 && langs[0].Tag != "" {
		locale = langs[0].Tag
	}

	main := protected.ProtectedMain(tr, protected.ProtectedParams{
		User:     userMeta,
		Attrs:    sessionAttrs,
		Data:     sessionData,
		Sessions: sessions,
		Active:   r.URL.Path,
		Locale:   locale,
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
