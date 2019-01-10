package web

import (
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"

	"github.com/danikarik/yandex-audience/pkg/audience"
	"github.com/danikarik/yandex-audience/pkg/config"
	"github.com/go-chi/chi"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
)

// Web contains web pages routes.
type Web struct {
	api      string
	store    sessions.Store
	conf     config.Specification
	handler  http.Handler
	renderer *render.Render
	oauth    *oauth2.Config
	uploader *audience.Uploader
}

// New creates new instance of web container.
func New(store sessions.Store, conf config.Specification) *Web {
	web := &Web{
		api:      "https://api-audience.yandex.ru/v1/management/segments/upload_csv_file",
		store:    store,
		conf:     conf,
		oauth:    newOAuthConfig(conf),
		uploader: audience.New(),
	}

	rnd := render.New(render.Options{
		Layout:     "layouts/_base",
		Extensions: []string{".tmpl", ".html"},
	})
	web.renderer = rnd

	mux := chi.NewRouter()
	mux.Get("/", web.indexHandler)
	mux.Get("/callback", web.callbackHandler)
	mux.Get("/authorize", web.authorizeHandler)
	mux.Group(func(r chi.Router) {
		r.Use(web.LoginRequired, web.AuthExpired)
		r.Get("/upload", web.uploadHandler)
		r.Post("/process", web.processHandler)
		r.Post("/logout", web.logoutHandler)
	})
	web.handler = csrf.Protect([]byte(conf.SecretKey), csrf.Secure(!conf.Debug))(mux)

	return web
}

func newOAuthConfig(conf config.Specification) *oauth2.Config {
	yandexOAuthConfig := &oauth2.Config{
		RedirectURL:  conf.SiteURL + "/callback",
		ClientID:     conf.Yandex.CliendID,
		ClientSecret: conf.Yandex.ClientSecret,
		Endpoint:     yandex.Endpoint,
	}
	return yandexOAuthConfig
}

// Handler returns container's http handler.
func (w *Web) Handler() http.Handler {
	return w.handler
}
