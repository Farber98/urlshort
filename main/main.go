package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"urlshort"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	yamlFile := flag.String("yaml", "config.yaml", ".yaml filename to read from")
	jsonFile := flag.String("json", "config.json", ".json filename to read from")
	flag.Parse()

	yaml, err := os.ReadFile(*yamlFile)
	if err != nil {
		log.Printf("yamlFile parsing err   #%v ", err)
	}

	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}

	json, err := os.ReadFile(*jsonFile)
	if err != nil {
		log.Printf("jsonFile parsing err   #%v ", err)
	}

	jsonHandler, err := urlshort.JSONHandler(json, yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
