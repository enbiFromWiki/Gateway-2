package app

import (
	"encoding/json"
	"fmt"
	"gateway/backend/eventstream"
	"gateway/backend/wiki"
	"net/http"
	"strings"
)

type Service struct {
	mux    *http.ServeMux
	client *wiki.Client
}

func createDeps() Service {
	return Service{
		mux:    http.NewServeMux(),
		client: wiki.New("", "https://en.wikipedia.org/w/api.php", "json", 2),
	}
}

func Run() {
	eventstream.StartWMStream()
	service := createDeps()
	fmt.Println(service.client.GetLocalEditCounts([]string{"enbi", "Alt of enbi"}))
	initRoutes(service)
	http.ListenAndServe("127.0.0.1:8080", service.mux)
}

func (service *Service) handleContribs(w http.ResponseWriter, r *http.Request) {
	user := strings.TrimPrefix(r.URL.Path, "/api/contribs/")
	if user == "" {
		http.Error(w, "Missing user", 400)
	}

	client := service.client
	data, err := client.GetContribs(user)
	if err != nil {
		http.Error(w, "Bad upstream response", 502)
	}

	w.Write(data)
}

func (service *Service) handleEditCount(w http.ResponseWriter, r *http.Request) {
	userstring := strings.TrimPrefix(r.URL.Path, "/api/editcount/")
	users := strings.Split(userstring, "|")
	if len(users) == 0 {
		http.Error(w, "Missing user", 400)
	}

	client := service.client
	data, err := client.GetLocalEditCounts(users)
	if err != nil {
		http.Error(w, "Bad upstream response", 502)
	}

	output, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Bad upstream response", 502)
	}

	w.Write(output)
}

func initRoutes(service Service) {
	service.mux.HandleFunc("/api/contribs/", service.handleContribs)
	service.mux.HandleFunc("/api/editcount/", service.handleEditCount)
}
