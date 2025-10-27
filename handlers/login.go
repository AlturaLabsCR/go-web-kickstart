package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"app/config"
	"app/internal/db"
	"app/templates"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	content := templates.Login(h.Translator(r))
	templates.Base(content).Render(ctx, w)
}

func (h *Handler) SendVerification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "missing email", http.StatusBadRequest)
		return
	}

	// validate email format + rate limit (IP or email)

	queries := db.New(h.DB())
	tempKey, _ := generateTempKey()
	expires := time.Now().Add(5 * time.Minute).Unix()

	if _, err := queries.GetTempKey(ctx, email); err != nil {
		err := queries.InsertTempKey(ctx, db.InsertTempKeyParams{
			TempKeyEmail:       email,
			TempKey:            tempKey,
			TempKeyExpiresUnix: expires,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Log().Error("insert_temp_key", "error", err)
			return
		}
	} else {
		err := queries.UpdateTempKey(ctx, db.UpdateTempKeyParams{
			TempKey:            tempKey,
			TempKeyExpiresUnix: expires,
			TempKeyEmail:       email,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Log().Error("insert_temp_key", "error", err)
			return
		}
	}

	verifyURL := fmt.Sprintf("%s%sverify?key=%s&email=%s",
		hostURL(r),
		config.RootPrefix,
		url.QueryEscape(tempKey),
		url.QueryEscape(email),
	)

	body := fmt.Sprintf(
		"Confirma tu email:\n\n%s\n\nEste enlace expira en 5 minutos.",
		verifyURL,
	)

	if err := h.SMTPClient.SendText(
		config.ServerSMTPUser,
		[]string{email},
		"Email de confirmación",
		body,
	); err != nil {
		h.Log().Error(
			"verification_email_send_failed",
			"smtp_from", config.ServerSMTPUser,
			"smtp_to", email,
			"error", err,
		)
		templates.LoginError(h.Translator(r)).Render(ctx, w)
		return
	}

	h.Log().Info(
		"verification_email_sent",
		"smtp_from", config.ServerSMTPUser,
		"smtp_to", email,
	)

	templates.LoginVerify(h.Translator(r), email).Render(ctx, w)
}

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tempKey := r.URL.Query().Get("key")
	email := r.URL.Query().Get("email")

	if tempKey == "" || email == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	queries := db.New(h.DB())
	tempKeyDB, err := queries.GetTempKey(ctx, email)
	if err != nil {
		http.Error(w, "invalid or expired key", http.StatusBadRequest)
		return
	}

	if tempKeyDB.TempKey != tempKey {
		http.Error(w, "invalid key", http.StatusBadRequest)
		return
	}
	if time.Now().Unix() > tempKeyDB.TempKeyExpiresUnix {
		http.Error(w, "key expired", http.StatusBadRequest)
		return
	}

	queries.SetTempKeyUsed(ctx, email)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    tempKey,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   r.TLS != nil,
	})

	h.Log().Info("verification_success", "email", email)

	fmt.Fprintln(w, "Verificación completada")
}

// FIXME: This is NOT secure, use JWTs in production for secure sessions
func generateTempKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func hostURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}
