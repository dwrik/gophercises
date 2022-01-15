package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/dwrik/urlshortener/urlshort"
)

func main() {
	yamlFilename := flag.String("yaml", "redirect.yaml", "a yaml file with url mappings")
	jsonFilename := flag.String("json", "redirect.json", "a json file with url mappings")

	flag.Parse()

	var jsonPresent bool
	var yamlPresent bool

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "json":
			jsonPresent = true
		case "yaml":
			yamlPresent = true
		}
	})

	var mapHandler, jsonHandler, yamlHandler http.HandlerFunc
	var jsonPaths, yamlPaths = []byte{}, []byte{}

	if jsonPresent {
		json, err := readFile(*jsonFilename)
		if err != nil {
			exit(fmt.Sprintf("failed to read json file: %s", *jsonFilename))
		}
		jsonPaths = json
	}

	if yamlPresent {
		yaml, err := readFile(*yamlFilename)
		if err != nil {
			exit(fmt.Sprintf("failed to read yaml file: %s", *yamlFilename))
		}
		yamlPaths = yaml
	}

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mux := defaultMux()
	mapHandler = urlshort.MapHandler(pathsToUrls, mux)

	yamlHandler, err := urlshort.YAMLHandler(yamlPaths, mapHandler)
	if err != nil {
		exit(fmt.Sprintf("invalid yaml file: %s", *yamlFilename))
	}

	jsonHandler, err = urlshort.JSONHandler(jsonPaths, yamlHandler)
	if err != nil {
		exit(fmt.Sprintf("invalid json file: %s", *jsonFilename))
	}

	fmt.Println("Starting server on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", jsonHandler))
}

func readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	read, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return read, nil
}

func exit(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello there!")
}
