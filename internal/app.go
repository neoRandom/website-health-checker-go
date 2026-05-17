package internal

import (
	"http-server/template"
	"net/http"
	"time"
	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
)

//

type GlobalState struct {
	Count int
}

var global GlobalState
var sessionManager *scs.SessionManager

//

type App struct {
	Addr              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func (a *App) GetServer() *http.Server {
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()

			// Check to see if the global button was pressed.
			if r.Form.Has("global") {
				global.Count++
			}
			if r.Form.Has("session") {
				currentCount := sessionManager.GetInt(r.Context(), "count")
				sessionManager.Put(r.Context(), "count", currentCount+1)
			}
		}

		sessionCount := sessionManager.GetInt(r.Context(), "count")
		component := template.Index(global.Count, sessionCount)
		component.Render(r.Context(), w)
	}))

	mux.Handle("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.FormValue("name")
		component := template.Hello(name)
		component.Render(r.Context(), w)
	}))

	mux.Handle("/404", templ.Handler(template.NotFound(), templ.WithStatus(http.StatusNotFound)))

	sMux := sessionManager.LoadAndSave(mux)
	wMux := WrapAllMiddleware(sMux)

	srv := &http.Server{
		Addr:              a.Addr,
		Handler:           wMux,
		ReadTimeout:       a.ReadTimeout,
		ReadHeaderTimeout: a.ReadHeaderTimeout,
		WriteTimeout:      a.WriteTimeout,
		IdleTimeout:       a.IdleTimeout,
	}

	return srv
}
