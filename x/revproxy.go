package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	u, err := url.Parse("http://localhost:80")
	if err != nil {
		log.Fatal(err)
	}
	p := httputil.NewSingleHostReverseProxy(u)
	log.Fatal(http.ListenAndServe(":9001", p))
}
