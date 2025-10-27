package middleware

import "net/http"

func Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// verify cookie validity (e.g. check against DB, expiration)
		// validate that this cookie corresponds to a valid session.

		// Continue to next handler if authorized
		next.ServeHTTP(w, r)
	})
}
