package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

var ErrCacheMiss = errors.New("cache miss")

// CachingRoundTripper keeps an in memory cache of response bodies.
type CachingRoundTripper struct {
	cache map[string][]byte
	rt    http.RoundTripper
}

func New() *CachingRoundTripper {
	return &CachingRoundTripper{
		cache: make(map[string][]byte),
		rt:    http.DefaultTransport,
	}
}

func (c *CachingRoundTripper) Set(key string, value []byte) {
	c.cache[key] = value
}

func (c *CachingRoundTripper) Get(key string) ([]byte, error) {
	if v, ok := c.cache[key]; ok {
		return v, nil
	}
	return nil, ErrCacheMiss
}

func (c *CachingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	key := req.URL.String()
	_, err := c.Get(key)
	if err != nil {
		log.Println("cache miss")
		resp, err := c.rt.RoundTrip(req)
		if err != nil {
			return resp, err
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}
		c.cache[key] = b
	} else {
		log.Printf("cache hit for %s", key)
	}
	return &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader(c.cache[key])),
	}, nil
}

func main() {
	client := http.Client{
		Transport: New(),
	}
	var err error

	for i := 0; i < 3; i++ {
		_, err = client.Get("http://golangleipzig.space")
		if err != nil {
			log.Fatal(err)
		}
	}
}
