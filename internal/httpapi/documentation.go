package httpapi

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/saeedAdelpour/paymentshare/api"
)

//go:embed docs/*
var documentation embed.FS

func registerDocumentation(router chi.Router) {
	router.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = w.Write(api.Document)
	})
	router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})
	// The embedded subtree always exists; panic would indicate a build-time mistake.
	files, err := fs.Sub(documentation, "docs")
	if err != nil {
		panic(err)
	}
	static := http.StripPrefix("/docs/", http.FileServer(http.FS(files)))
	router.Get("/docs/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		static.ServeHTTP(w, r)
	})
}
