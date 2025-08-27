package utils

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"net/http"
)

const client_id string = "593200518603-k0ptna6taq593eiulqnd4vfsk1djh0vl.apps.googleusercontent.com"

// https://developers.google.com/identity/protocols/oauth2/#installed
const client_nonsecret string = "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"

func LogIn() error {
	return nil
}

// See: https://developers.google.com/identity/protocols/oauth2/native-app#step1-code-verifier
func pkce() (string, string, error) {
	codeVerifier, err := generateRandomBytes(64)
	if err != nil {
		return "", "", err
	}
	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(hash[:])

	return codeVerifier, codeChallenge, nil
}

func generateRandomBytes(len int) (string, error) {
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil || len < 0 {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func start() error {
	config := &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_nonsecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{drive.DriveFileScope},
		RedirectURL:  "http://localhost:8080/callback",
	}

	code_verifier, code_challenge, err := pkce()

	if err != nil {
		return err
	}

	state, err := generateRandomBytes(16)
	if err != nil {
		return err
	}

	authURL := config.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", code_challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))

	// call browser
	fmt.Println(authURL)

	codeCh := make(chan string)
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		fmt.Fprintln(w, "Auth complete! You can close this window.")
		codeCh <- code
	})
	go http.ListenAndServe(":8080", nil)
	code := <-codeCh

	// Exchange with PKCE
	token, err := config.Exchange(context.Background(), code,
		oauth2.SetAuthURLParam("code_verifier", code_verifier),
	)
	if err != nil {
		fmt.Println("error:", err.Error())
		return err
	}

	fmt.Println("Access Token:", token.AccessToken)

	return nil
}
