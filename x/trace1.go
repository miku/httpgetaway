package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
)

func main() {
	req, err := http.NewRequest("GET", "https://golangleipzig.space", nil)
	if err != nil {
		log.Fatal(err)
	}
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", dnsInfo)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	_, err = http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Fatal(err)
	}
}
