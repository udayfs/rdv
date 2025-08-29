package utils

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	// https://developers.google.com/identity/protocols/oauth2/#installed
	clientNonsecret string = "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"
	clientID        string = "593200518603-k0ptna6taq593eiulqnd4vfsk1djh0vl.apps.googleusercontent.com"
	port            string = "8008"
	tokenFile       string = ".rdv_user.json"
)

var (
	redirectURL string = fmt.Sprintf("http://localhost:%s/callback", port)
	userHome, _        = os.UserHomeDir()
	oauthConfig        = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientNonsecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{drive.DriveFileScope},
		RedirectURL:  redirectURL,
	}
)

func generateRandomBytes(len int) (string, error) {
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil || len < 0 {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
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

func buildToken(tokenData map[string]string) (*oauth2.Token, error) {
	expiry, err := time.Parse("2006-01-02T15:04:05Z07:00", tokenData["expiry"])
	if err != nil {
		return nil, err
	}
	return &oauth2.Token{
		AccessToken:  tokenData["access_token"],
		RefreshToken: tokenData["refresh_token"],
		TokenType:    tokenData["token_type"],
		Expiry:       expiry,
	}, nil
}

func saveToken(token *oauth2.Token) error {
	if userHome == "" {
		return os.ErrNotExist
	}

	tokenJSON := filepath.Join(userHome, tokenFile)
	tokenData := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"token_type":    token.TokenType,
		"expiry":        token.Expiry.Format("2006-01-02T15:04:05Z07:00"),
	}

	file, err := os.Create(tokenJSON)
	if err != nil {
		return err
	}

	jsonEncoder := json.NewEncoder(file)
	jsonEncoder.SetIndent("", "  ")
	if err := jsonEncoder.Encode(tokenData); err != nil {
		return err
	}

	return file.Close()
}

func authorize() error {
	code_verifier, code_challenge, err := pkce()

	if err != nil {
		return err
	}

	state, err := generateRandomBytes(16)
	if err != nil {
		return err
	}

	authURL := oauthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", code_challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))

	if err := OpenURL(authURL); err != nil {
		return err
	}

	codeChn := make(chan string)
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code := query.Get("code")
		rState := query.Get("state")

		if rState != state {
			http.Error(w, "Invalid 'state' parameter!", http.StatusBadRequest)
			codeChn <- ""
			return
		}

		fmt.Fprintln(w, "Authorization successful! You can close this window now.")
		codeChn <- code
	})

	go func() {
		http.ListenAndServe(":"+port, nil)
	}()
	code := <-codeChn

	// exchange with pkce
	token, err := oauthConfig.Exchange(context.Background(), code,
		oauth2.SetAuthURLParam("code_verifier", code_verifier),
	)

	if err != nil {
		return err
	}

	return saveToken(token)
}

func LogIn() (*http.Client, error) {
	tokenData := make(map[string]string)
	if userHome == "" {
		return nil, os.ErrNotExist
	}

	file, err := os.Open(filepath.Join(userHome, tokenFile))
	if err == nil {
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&tokenData); err == nil {
			if token, err := buildToken(tokenData); err == nil {
				client := oauthConfig.Client(context.Background(), token)
				return client, nil
			}
		}
	}

	if err := authorize(); err != nil {
		return nil, err
	}

	file, err = os.Open(filepath.Join(userHome, tokenFile))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&tokenData); err != nil {
		return nil, err
	}

	token, err := buildToken(tokenData)
	if err != nil {
		return nil, err
	}

	client := oauthConfig.Client(context.Background(), token)
	return client, nil
}
