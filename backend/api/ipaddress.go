package api

import (
	"net"
	"net/http"
)

func GetIp(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")

	if ip != "" {
		return ip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return r.RemoteAddr
	}

	return ip
}