package auth

import (
	"context"
	"net/http"

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
	ctx := context.WithValue(
		context.Background(),
		oauth2.HTTPClient,
		&http.Client{},
	)
	code := c.Query("code")

	if code == "" {
		c.String(400, "No oauth2 code returned")
		return
	}

	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		c.String(500, "Token exchange failed: %t", err)
		return
	}

	c.JSON(200, gin.H{
		"access_token": token.AccessToken,
		"token_type":   token.TokenType,
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
