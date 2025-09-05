package handler

import (
	"net/http"

	"github.com/Samuel-Martineau/watsearch/internal/templates"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Mux() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// panic("Hello")
		templates.Hello("Samuel").Render(r.Context(), w)
	})

	r.Get("/path", templ.Handler(templates.Hello("Dang Khoa")).ServeHTTP)

	return r
}
