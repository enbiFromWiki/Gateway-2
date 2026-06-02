package main

import (
	"encoding/json"
	"fmt"
	"gateway/backend/api"
	"gateway/backend/mw"
	"net/http"
	"time"
)

type SiteInfoResponse struct {
	Query struct {
		General struct {
			SiteName   string `json:"sitename"`
			MainPage   string `json:"mainpage"`
			Lang       string `json:"lang"`
			Generator  string `json:"generator"`
			ServerName string `json:"servername"`
			WikiID     string `json:"wikiid"`
		} `json:"general"`
	} `json:"query"`
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")

	clientIp := api.GetIp(r)
	w.Write([]byte(clientIp))
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request: %s %s @ %s\n", r.Method, r.URL.Path, time.Now().Format("15:04:05"))

		next.ServeHTTP(w, r)
	})
}

func main() {
	client := mw.Client{
		HTTPclient: &http.Client{},
		DefaultUrl: "https://en.wikipedia.org/w/api.php",
		User:       "enbi",
	}

	data, err := client.Get(map[string]string{
		"action": "query",
		"meta":   "siteinfo",
		"siprop": "general",
		"format": "json",
	})
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}

	var info SiteInfoResponse

	err = json.Unmarshal(data, &info)

	fmt.Println(info.Query.General.SiteName)
	fmt.Println(info.Query.General.MainPage)
	fmt.Println(info.Query.General.Generator)

}
