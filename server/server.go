package server 

import (
	"fmt"
	"net/http"
	"os"
  "strconv"
  "time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
  port int
  spotifyClient *SpotifyClient
}

func NewServer() *http.Server {
  port, _ := strconv.Atoi(get_port())

  NewServer := &Server {
    port: port,
    spotifyClient: NewClient(),
  };

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server

}


type User struct {
	UID               string
	Slug              string
	AuthorizationCode string
  // Settings scalar?
}

func get_port() string {
  port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}


