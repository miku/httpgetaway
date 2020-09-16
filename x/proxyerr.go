package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

func main() {
	var transport http.RoundTripper = &http.Transport{
		Proxy: func(r *http.Request) (*url.URL, error) {
			return nil, fmt.Errorf("proxy err")
		},
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := http.Client{
		Transport: transport,
	}
	resp, err := client.Get("https://heise.de")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", resp)
}
