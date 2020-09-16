package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	flag.Parse()
	url := flag.Arg(0)
	if url == "" {
		url = "http://heise.de"
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", resp)
}

// HTTP_PROXY=http://fluxproxy.com:1080 go run clientproxy.go https://heise.de
// 2020/09/16 11:53:48 Get "https://heise.de": proxyconnect tcp: dial tcp 159.69.240.245:1080: connect: connection refused
