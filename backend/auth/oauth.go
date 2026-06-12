package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var oauthConfig = &oauth2.Config{
	ClientID:     "a1ccc2fee4803887ae7f39acee5ed81d",
	ClientSecret: "53b966b3aa43726e0a5ea7d05627a396b1c8e990",
	RedirectURL:  "http://localhost:8080/auth/callback",
	Scopes: []string{
		"basic",
		"edit", // not a clue what should be here.
	},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://meta.wikimedia.org/w/rest.php/oauth2/authorize",
		TokenURL: "https://meta.wikimedia.org/w/rest.php/oauth2/access_token",
	},
}

func Login(c *gin.Context) {
	url := oauthConfig.AuthCodeURL("random-state-string")
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
