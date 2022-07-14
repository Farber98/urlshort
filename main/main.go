package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"urlshort"

	"github.com/boltdb/bolt"
)

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = writeDb(db)
	if err != nil {
		log.Fatal(err)
	}

	pathsToUrls, err := readDb(db)
	if err != nil {
		log.Fatal(err)
	}

	mux := defaultMux()

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
