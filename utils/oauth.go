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

type OauthProvider struct {
	name     string
	clientID string
	// https://developers.google.com/identity/protocols/oauth2/#installed
	clientSec string
	scopes    []string
	endpoint  oauth2.Endpoint
}

var (
	globalConfig  *oauth2.Config
	port          string = "8008"
	redirectURL   string = fmt.Sprintf("http://localhost:%s/callback", port)
	tokenFileName string = ".rdv.json"
	tokenFilePath string = filepath.Join(getUserHome(), tokenFileName)
)

var providers = []OauthProvider{
	{
		name:      "google",
		clientID:  "593200518603-k0ptna6taq593eiulqnd4vfsk1djh0vl.apps.googleusercontent.com",
		clientSec: "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",
		endpoint:  google.Endpoint,
		scopes:    []string{drive.DriveFileScope},
	},
}

func getUserHome() string {
	res, err := os.UserHomeDir()
	if err != nil {
		ExitOnError(err.Error())
	}
	return res
}

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

func getToken(tokenData map[string]string) (*oauth2.Token, error) {
	file, err := os.Open(tokenFilePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	if err := json.NewDecoder(file).Decode(&tokenData); err != nil {
		return nil, err
	}

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

func setToken(token *oauth2.Token, tokenFile string) error {
	tokenData := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"token_type":    token.TokenType,
		"expiry":        token.Expiry.Format("2006-01-02T15:04:05Z07:00"),
	}

	file, err := os.Create(tokenFile)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonEncoder := json.NewEncoder(file)
	jsonEncoder.SetIndent("", "  ")

	return jsonEncoder.Encode(tokenData)
}

func authorize(provider *OauthProvider) (*http.Client, error) {
	switch provider.name {
	case "google":
	default:
		return nil, fmt.Errorf("unsupported oauth provider")
	}

	var oauthConfig = &oauth2.Config{
		ClientID:     provider.clientID,
		ClientSecret: provider.clientSec,
		Endpoint:     provider.endpoint,
		RedirectURL:  redirectURL,
		Scopes:       provider.scopes,
	}

	globalConfig = oauthConfig

	code_verifier, code_challenge, err := pkce()

	if err != nil {
		return nil, err
	}

	state, err := generateRandomBytes(16)
	if err != nil {
		return nil, err
	}

	authURL := oauthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", code_challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))

	if err := OpenURL(authURL); err != nil {
		return nil, err
	}

	codeChn := make(chan string)
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code := query.Get("code")
		rState := query.Get("state")
		defer close(codeChn)

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
		return nil, err
	}

	if err := setToken(token, tokenFilePath); err != nil {
		return nil, err
	}

	client := oauthConfig.Client(context.Background(), token)
	return client, nil
}

func LogIn() (*http.Client, error) {
	tokenData := make(map[string]string)

	if token, err := getToken(tokenData); err == nil {
		return globalConfig.Client(context.Background(), token), nil
	}

	// only google for now
	client, err := authorize(&providers[0])
	if err != nil {
		return nil, err
	}

	return client, err
}
