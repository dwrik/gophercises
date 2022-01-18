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
	yamlFilename := flag.String("yaml", "", "a yaml file with url mappings")
	jsonFilename := flag.String("json", "", "a json file with url mappings")

	flag.Usage = customUsage()
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

	var handler http.Handler

	switch {
	case jsonPresent:
		handler = getHandler(urlshort.JSONHandler, *jsonFilename)
	case yamlPresent:
		handler = getHandler(urlshort.YAMLHandler, *yamlFilename)
	default:
		exit("Error:\tConfiguration file with url mappings not provided.")
	}

	fmt.Println("Starting server on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func customUsage() func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Usage: urlshortener OPTION FILE\n\n")
		flagNames := [...]string{"json", "yaml"}
		for _, fname := range flagNames {
			f := flag.Lookup(fname)
			fmt.Fprintf(os.Stderr, "  -%s\t%s\n", f.Name, f.Usage)
		}
		fmt.Fprintln(os.Stderr, "\nIf both -json and -yaml files are\nprovided then the -json FILE is used.")
	}
}

func getHandler(handlerFactory func([]byte, http.Handler) (http.HandlerFunc, error), filename string) http.Handler {
	pathsToUrls := map[string]string{
		"/yaml-godoc": "https://godoc.org/gopkg.in/yaml.v2",
	}

	mux := defaultMux()
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	file, err := readFile(filename)
	if err != nil {
		exit(fmt.Sprintf("failed to read file:\t%s", filename))
	}

	handler, err := handlerFactory(file, mapHandler)
	if err != nil {
		exit(fmt.Sprintf("invalid file:\t%s", filename))
	}

	return handler
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
