package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type User struct {
	UID               string
	Slug              string
	AuthorizationCode string
}

type CurrentlyPlaying struct {
	Item struct {
		Album struct {
			Name string `json:"name"`
		} `json:"album"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
		Name string `json:"name"`
	} `json:"item"`
}

const (
	tokenURL         = "https://accounts.spotify.com/api/token"
	redirectURI      = "http://localhost:8080/auth/callback/spotify"
	spotifyAppScopes = "user-read-currently-playing"
)

var (
	clientID     = ""
	clientSecret = ""
)

func get_port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID = os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")

	port := get_port()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/auth/callback/spotify", spotifyCallback)

	http.HandleFunc("/now-playing", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", "to play a song")
	})

	fmt.Printf("Starting server at port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomString(16)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r,
		"https://accounts.spotify.com/authorize?client_id="+clientID+
			"&response_type=code&redirect_uri="+redirectURI+"&scope="+spotifyAppScopes+
			"&state="+state,
		http.StatusFound)
}

func spotifyCallback(w http.ResponseWriter, r *http.Request) {
	// state := r.FormValue("state")
	// check state

	code := r.FormValue("code")
	token, err := getAccessToken(code)
	if err != nil {
		fmt.Printf("Error exchanging code for token: %v\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	currentlyPlaying, err := getNowPlaying(token.AccessToken)
	if err != nil {
		fmt.Printf("Error getting currently playing song: %v\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Printf("Currently playing: %v\n", currentlyPlaying)

	// Store the users token
	// User {
	// UID string/int/serial/something
	// url for now playing
	// token
	// settings
}

func handleNowPlaying(w http.ResponseWriter, r *http.Request) {
	_, err := getNowPlaying("")
	if err != nil {
		fmt.Printf("Error getting currently playing song: %v\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
}

func getNowPlaying(token string) (*CurrentlyPlaying, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)

	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("Error getting currently playing song: %v\n", err)
		return nil, err
	}

	defer resp.Body.Close()

	currentlyPlaying := &CurrentlyPlaying{}

	if err := json.NewDecoder(resp.Body).Decode(currentlyPlaying); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return nil, err
	}

	return currentlyPlaying, nil
}

func getAccessToken(code string) (*Token, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)
	form.Set("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	form.Set("client_secret", os.Getenv("SPOTIFY_CLIENT_SECRET"))

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	token := &Token{}
	if err := json.NewDecoder(resp.Body).Decode(token); err != nil {
		return nil, err
	}

	return token, nil

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
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

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
