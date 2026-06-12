package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"crypto/rand"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

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
		AuthURL:  "https://meta.wikimedia.org/wiki/Special:OAuth/approve",
		TokenURL: "https://meta.wikimedia.org/w/rest.php/oauth2/access_token",
	},
}

// func Login (*gin.Context) {
// 	url, _ := url.Parse("https://meta.wikimedia.org/wiki/Special:OAuth/approve")
// 	query := url.Query()
// 	query.Set("returnto", "http://localhost:8080")

// }

func generateRandomCode() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	output := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)

	return output, nil
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

func Callback(c *gin.Context) {
	// ctx := context.WithValue(
	// 	context.Background(),
	// 	oauth2.HTTPClient,
	// 	&http.Client{},
	// )
	code := c.Query("code")

	if code == "" {
		c.String(400, "No oauth2 code returned")
		return
	}

	// token, err := oauthConfig.Exchange(ctx, code) lol this just doesnt work

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", "http://localhost:8080/auth/callback")

	req, _ := http.NewRequest("POST",
		"https://meta.wikimedia.org/w/rest.php/oauth2/access_token",
		strings.NewReader(data.Encode()),
	)

	req.SetBasicAuth(oauthConfig.ClientID, oauthConfig.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		c.String(500, "Token exchange failed: %t", err.Error())
		return
	}

	defer res.Body.Close()

	var result map[string]any

	json.NewDecoder(res.Body).Decode(&result)

	c.JSON(200, result)
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
