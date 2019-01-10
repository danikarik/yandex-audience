package web

import (
	"net/http"

	"github.com/gorilla/csrf"
)

type htmlData struct {
	IsAuthenticated bool
	CsrfToken       string
	Flashes         []string
}

func (web *Web) token(r *http.Request) (string, bool) {
	sess, _ := web.store.Get(r, "_ya_session")
	token, ok := sess.Values["access_token"].(string)
	return token, ok
}

func (web *Web) isAuthenticated(r *http.Request) bool {
	token, ok := web.token(r)
	return token != "" && ok
}

func (web *Web) flashes(r *http.Request) []string {
	sess, _ := web.store.Get(r, "_ya_flashes")
	f := make([]string, 0)
	for _, m := range sess.Flashes("message") {
		s, ok := m.(string)
		if ok {
			f = append(f, s)
		}
	}
	return f
}

func (web *Web) html(w http.ResponseWriter, r *http.Request, name string) {
	htmlData := &htmlData{
		IsAuthenticated: web.isAuthenticated(r),
		CsrfToken:       csrf.Token(r),
		Flashes:         web.flashes(r),
	}
	web.renderer.HTML(w, http.StatusOK, name, htmlData)
}
