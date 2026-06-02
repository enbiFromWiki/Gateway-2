package mw

import "net/http"

var Wikipedia = &Client{
	HTTPclient: &http.Client{},
	DefaultUrl: "https://en.wikipedia.org/w/api.php",
	User:       "enbi",
	UserAgent:  "User:enbi's test script in Go",
}
