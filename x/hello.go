package main

import (
	"log"
	"net/http"
)

func main() {
	loc := "https://golangleipzig.space"
	resp, err := http.Get(loc)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Printf("%d %s", resp.StatusCode, loc)
}
