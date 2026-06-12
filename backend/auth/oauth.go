package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"crypto/rand"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type MWOauth struct {
	config *oauth2.Config
	ua     string
}

var oauthConfig = &oauth2.Config{
	ClientID:     "074bc5c055a61844e6fdd4f91d7ef345",
	ClientSecret: "fa7f86af755ccdf76aade8c6f7b953dad4bc2e79",
	RedirectURL:  "http://localhost:8080/auth/callback",
	Scopes: []string{
		"basic",
		"editpage",
		"rollback",
	},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://meta.wikimedia.org/w/rest.php/oauth2/authorize",
		TokenURL: "https://meta.wikimedia.org/w/rest.php/oauth2/access_token",
	},
}

// func Login (*gin.Context) {
// 	url, _ := url.Parse("https://meta.wikimedia.org/wiki/Special:OAuth/approve")
// 	query := url.Query()
// 	query.Set("returnto", "http://localhost:8080")

// }

// func ()

func generateRandomCode() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	output := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)

	return output, nil
}

var authenticator = &MWOauth{
	config: oauthConfig,
	ua:     "User:enbi/OAuth Testing (localhost dev)",
}

func Login(c *gin.Context) {
	state, err := generateRandomCode()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating random string: %t", err)
		return
	}

	url := oauthConfig.AuthCodeURL(state) + "&oauth_version=2"
	fmt.Println(url)
	c.Redirect(302, url)
}

func (t *MWOauth) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.ua)
	return http.DefaultTransport.RoundTrip(req)
}

func (a *MWOauth) getToken(code string) (*oauth2.Token, error) {
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: a.config.TokenSource(context.Background(), &oauth2.Token{
				AccessToken: code,
			}),
			Base: a,
		},
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)

	return a.config.Exchange(ctx, code)
}

func Callback(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		c.String(400, "No oauth2 code returned")
		return
	}

	token, err := authenticator.getToken(code)

	if err != nil {
		c.String(500, "Token exchange failed: %t", err.Error())
		return
	}

	fmt.Println("TOKEN::: " + token.AccessToken)

	c.JSON(200, gin.H{
		"status": "success",
		"token":  token.AccessToken,
	})
}

func ApiTest(c *gin.Context) {
	token := c.Query("token")
	client := oauthConfig.Client(context.Background(), &oauth2.Token{
		AccessToken: token,
	})

	res, err := client.Get("https://test/wikipedia.org/w/api.php?action=query&meta=tokens")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to get CSRF token: %t", err)
	}

	c.JSON(200, gin.H{
		"status":   "success",
		"response": res,
	})
}
