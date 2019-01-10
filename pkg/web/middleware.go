package web

import (
	"net/http"
	"time"
)

// LoginRequired checks if access_token present.
func (web *Web) LoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, _ := web.store.Get(r, "_ya_session")
		if token, ok := sess.Values["access_token"].(string); !ok || token == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuthExpired checks if token is not expired.
func (web *Web) AuthExpired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, _ := web.store.Get(r, "_ya_session")
		if expiry, ok := sess.Values["expiry"].(time.Time); !ok || time.Now().After(expiry) {
			sess.Options.MaxAge = -1
			err := sess.Save(r, w)
			if err != nil {
				http.Error(w, http.StatusText(401), 401)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
