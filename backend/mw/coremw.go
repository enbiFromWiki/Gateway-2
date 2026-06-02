package mw

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	HTTPclient    *http.Client
	UserAgent     string
	DefaultUrl    string
	Formatversion int8
	User          string
}

func (client *Client) Get(params map[string]string) ([]byte, error) {
	// params["assert"] = "user"
	// if user := client.User; len(user) != 0 {
	// 	params["assertuser"] = user
	// }

	fmt.Println("Request sent")

	url := client.DefaultUrl + "?" + MapToQueryParams(params)

	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}

	if agent := client.UserAgent; len(agent) != 0 {
		req.Header.Set("User-Agent", agent)
	} else {
		req.Header.Set("User-Agent", "User:enbi's test script")
	}

	fmt.Println(req.Header.Get("User-Agent"))

	req.Header.Set("User-Agent", "User:enbi's test script")

	resp, err := client.HTTPclient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"API returned %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return data, nil

}

// func New(userAgent string, defaultUrl *url.URL) *Client {

// }

func MapToQueryParams(params map[string]string) string {
	values := url.Values{}

	for k, v := range params {
		values.Set(k, v)
	}

	return values.Encode()
}
