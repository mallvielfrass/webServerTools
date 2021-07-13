package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mallvielfrass/wst"
)

func register(w http.ResponseWriter, r *http.Request) {

}
func login(w http.ResponseWriter, r *http.Request) {

}
func profile(w http.ResponseWriter, r *http.Request) {

}
func main() {
	//path := wst.StaticFolder{Path: "./static"}
	r := chi.NewRouter()
	r.Use(wst.MiddlewareAllowCORS)
	r.Use(wst.MiddlewareURL)
	r.With(wst.MiddlewareJSON).Route("/api", func(r chi.Router) {
		//only not auth methods
		r.With().Route("/nauth", func(r chi.Router) {
			r.HandleFunc("/register", register)
			r.HandleFunc("/login", login)
		})
		//only auth methods
		r.With().Route("/auth", func(r chi.Router) {
			r.HandleFunc("/profile", profile)
		})
	})
	wst.FileServer(r, "static")
	http.ListenAndServe(":3333", r)
}
