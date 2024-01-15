package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type SpotifyClient struct {
	clientID     string
	clientSecret string
	redirectURL  string
	scopes       string
	tokenURL     string
}

func NewClient() *SpotifyClient {
	return &SpotifyClient{
		clientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		clientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		redirectURL:  os.Getenv("SPOTIFY_REDIRECT_URL"),
		scopes:       "user-read-currently-playing",
		tokenURL:     "https://accounts.spotify.com/api/token",
	}
}

func (c *SpotifyClient) GetAuthURL(state string) string {
	return "https://accounts.spotify.com/authorize?client_id=" + c.clientID +
		"&response_type=code&redirect_uri=" + c.redirectURL + "&scope=" + c.scopes +
		"&state=" + state
}

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func (c *SpotifyClient) getAccessToken(code string) (*Token, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", c.redirectURL)
	form.Set("client_id", c.clientID)
	form.Set("client_secret", c.clientSecret)

	resp, err := http.Post(c.tokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
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

func (c *SpotifyClient) getNowPlaying(token string) (*CurrentlyPlaying, error) {
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
