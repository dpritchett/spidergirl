package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// HTTP GET a URL, return body as a string
func curl(uri string) ([]byte, error) {
	resp, err := http.Get(uri)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

// Curl text at a URL and send its response body to a handler channel
func curlAndLog(uri string, handler chan string) {
	log.Println("Curling ", uri)
	result, err := curl(uri)

	if err != nil {
		log.Fatal(err)
	}

	handler <- string(result)
}

// Curl a list of URLs concurrently and print out their results before exiting.
func main() {
	uris := []string{
		"http://localhost:3000/products?page=1",
		"http://localhost:3000/products?page=2",
		"http://localhost:3000/products?page=3",
		"http://localhost:3000/products?page=4",
	}

	results := make(chan string)

	// fire off a bunch of concurrent readers
	for _, uri := range uris {
		go curlAndLog(uri, results)
	}

	// Listen for results from the readers and print them to the screen
	for i := 0; i < len(uris); i++ {
		body := <-results
		fmt.Println(body[0:255])
	}
}
