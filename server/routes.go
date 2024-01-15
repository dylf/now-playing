package server

import (
	"log"
  "fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.homeHandler)
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

	http.Redirect(w, r, s.spotifyClient.GetAuthURL(state), http.StatusFound)
}

func (s *Server) spotifyCallback(w http.ResponseWriter, r *http.Request) {
	// state := r.FormValue("state")
	// check state

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

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	content := `
  <!DOCTYPE html>
  <html>
    <head>
      <title>Home Page</title>
  </head>
  <body style="background: black; color: white">
    <h1>Home Page</h1>
    <p>Welcome to my home page!</p>
    <a href="/login">Sign In</a>
  </body>
  </html>
  `
	w.Write([]byte(content))
}
