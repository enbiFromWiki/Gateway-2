package eventstream

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func StartWMStream() {
	req, err := http.NewRequest("GET", "https://stream.wikimedia.org/v2/stream/recentchange", nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	req.Header.Set("User-Agent", "User:enbi's test script in Go")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	scanner := bufio.NewScanner(res.Body)

	for scanner.Scan() {
		if e := scanner.Err(); e != nil {
			log.Fatal(e)
		}
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		fmt.Println(line)
	}
}
