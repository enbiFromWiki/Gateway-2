package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var oauthConfig = &oauth2.Config{
	ClientID:     "c873b63b13b3aafdbcebe29eb273e83d",
	ClientSecret: "d358df7e6a7712986cd499525d1f39d2d0485eef",
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

func Login(c *gin.Context) {
	url := oauthConfig.AuthCodeURL("random-state-string")
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
