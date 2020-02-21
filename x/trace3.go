package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"strings"
)

// transport is an http.RoundTripper that keeps track of the in-flight request
// and implements hooks to report HTTP tracing events. Also, keep resp around,
// so we can emit status code.
type transport struct {
	current *http.Request
}

// RoundTrip wraps http.DefaultTransport.RoundTrip to keep track
// of the current request.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.current = req
	return http.DefaultTransport.RoundTrip(req)
}

// GotConn prints whether the connection has been used previously
// for the current request.
func (t *transport) GotConn(info httptrace.GotConnInfo) {
	fmt.Printf("Connection reused for %v? %v\n", t.current.URL, info.Reused)
}

func prependHTTP(s string) string {
	if !strings.HasPrefix(s, "http") {
		return "http://" + s
	}
	return s
}

func main() {
	site := "https://golangleipzig.space"
	flag.Parse()
	if flag.NArg() > 0 {
		site = prependHTTP(flag.Arg(0))
	}
	t := &transport{}

	req, _ := http.NewRequest("GET", site, nil)
	trace := &httptrace.ClientTrace{
		GotConn: t.GotConn,
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	client := &http.Client{Transport: t}
	if _, err := client.Do(req); err != nil {
		log.Fatal(err)
	}
	if _, err := client.Do(req); err != nil {
		log.Fatal(err)
	}
}
