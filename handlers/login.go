package handlers

import (
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"time"

	"app/config"
	"app/internal/db"
	"app/templates"
	"app/utils"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	content := templates.Login(h.Translator(r))
	templates.Base(content).Render(ctx, w)
}

func (h *Handler) SendVerification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	email := r.FormValue("email")

	if _, err := mail.ParseAddress(email); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	queries := db.New(h.DB())

	tempKey, _ := utils.RandomString()
	tempKeyHash, _ := utils.Hash(tempKey)
	expires := time.Now().Add(5 * time.Minute).Unix()

	if _, err := queries.GetTempKey(ctx, email); err != nil {
		queries.InsertTempKey(ctx, db.InsertTempKeyParams{
			TempKeyEmail:       email,
			TempKeyHash:        tempKeyHash,
			TempKeyExpiresUnix: expires,
		})
	} else {
		queries.UpdateTempKey(ctx, db.UpdateTempKeyParams{
			TempKeyEmail:       email,
			TempKeyHash:        tempKeyHash,
			TempKeyExpiresUnix: expires,
		})
	}

	verifyURL := fmt.Sprintf("%s%sverify?key=%s&email=%s",
		utils.HostURL(r),
		config.RootPrefix,
		url.QueryEscape(tempKey),
		url.QueryEscape(email),
	)

	body := fmt.Sprintf(
		"Confirma tu email:\n\n%s\n\nEste enlace expira en 5 minutos.",
		verifyURL,
	)

	h.SMTPClient().SendText(
		config.ServerSMTPUser,
		[]string{email},
		"Email de confirmación",
		body,
	)

	templates.LoginVerify(h.Translator(r), email).Render(ctx, w)
}

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	email := r.URL.Query().Get("email")
	tempKey := r.URL.Query().Get("key")

	if tempKey == "" || email == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	queries := db.New(h.DB())

	tempKeyDB, err := queries.GetTempKey(ctx, email)
	if err != nil {
		http.Error(w, "invalid or expired key", http.StatusBadRequest)
		return
	}

	if !utils.HashCompare(tempKeyDB.TempKeyHash, tempKey) {
		http.Error(w, "invalid key", http.StatusBadRequest)
		return
	}

	if time.Now().Unix() > tempKeyDB.TempKeyExpiresUnix {
		http.Error(w, "key expired", http.StatusBadRequest)
		return
	}

	queries.SetTempKeyUsed(ctx, email)

	h.Sessions.JWTSet(w, r, email)

	h.Log().Info("verification_success", "email", email)

	fmt.Fprintln(w, "Verificación completada")
}
