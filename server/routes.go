package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.HomeHandler)
	r.Get("/login", s.LoginHandler)
	r.Get("/auth/callback/spotify", s.spotifyCallback)
	r.Get("/now-playing", s.HandleNowPlaying)

	return r
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomString(16)
	if err != nil {
		log.Fatal(err)
	}

	stateCookie := http.Cookie{
		Name:  "state",
		Path:  "/",
		Value: state,
	}

	http.SetCookie(w, &stateCookie)
	http.Redirect(w, r, s.spotifyClient.GetAuthURL(state), http.StatusFound)
}

func (s *Server) spotifyCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	// check state
	stateCookie, err := r.Cookie("state")
	if err != nil || stateCookie.Value != state {
		log.Printf("Invalid state: %s\n", stateCookie.Value)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := s.spotifyClient.getAccessToken(code)
	if err != nil {
		log.Printf("Error exchanging code for token: %v\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Store the users token
	// User {
	// UID string/int/serial/something
	// url for now playing
	// token
	// settings

	currentlyPlaying, err := s.spotifyClient.getNowPlaying(token.AccessToken)
	if err != nil {
		log.Printf("Error getting currently playing song: %v\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	log.Printf("Currently playing: %v\n", currentlyPlaying)

}

func (s *Server) HandleNowPlaying(w http.ResponseWriter, r *http.Request) {
	_, err := s.spotifyClient.getNowPlaying("")
	if err != nil {
		fmt.Printf("Error getting currently playing song: %v\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
}

func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	content := `
  <!DOCTYPE html>
  <html>
    <head>
      <title>Now Playing</title>
  </head>
  <body style="background: black; color: white">
    <h1>Now Playing</h1>
    <p>Welcome to my app!</p>
    <a href="/login">Sign In</a>
  </body>
  </html>
  `
	w.Write([]byte(content))
}
