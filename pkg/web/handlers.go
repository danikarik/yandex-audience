package web

import (
	"fmt"
	"net/http"

	"github.com/danikarik/yandex-audience/pkg/audience"

	"golang.org/x/oauth2"
)

func (web *Web) indexHandler(w http.ResponseWriter, r *http.Request) {
	web.html(w, r, "index")
}

func (web *Web) uploadHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := web.store.Get(r, "_ya_flashes")
	sess.Options.MaxAge = -1
	if err := sess.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	web.html(w, r, "upload")
}

func (web *Web) authorizeHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := web.store.Get(r, "_ya_session")
	opts := oauth2.SetAuthURLParam("response_type", "code")
	url := web.oauth.AuthCodeURL(sess.ID, opts)
	http.Redirect(w, r, url, http.StatusFound)
}

func (web *Web) callbackHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := web.store.Get(r, "_ya_session")
	code := r.FormValue("code")
	token, err := web.oauth.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if !token.Valid() {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	sess.Values["access_token"] = token.AccessToken
	sess.Values["expiry"] = token.Expiry
	if err = sess.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/upload", http.StatusSeeOther)
}

func (web *Web) logoutHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := web.store.Get(r, "_ya_session")
	sess.Options.MaxAge = -1
	if err := sess.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (web *Web) processHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := web.store.Get(r, "_ya_flashes")
	f, fh, err := r.FormFile("csvfile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, _ := web.token(r)
	p := &audience.Payload{}
	p, err = web.uploader.Do(token, fh.Filename, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sess.AddFlash(fmt.Sprintf("File successfully uploaded. Segment: %d", p.Segment.ID), "message")
	err = sess.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/upload", http.StatusSeeOther)
}
