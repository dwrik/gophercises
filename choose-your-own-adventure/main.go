package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dwrik/choose-your-own-adventure/cyoa"
)

var tmpl = template.Must(template.ParseFiles("templates/cyoa.html"))

type server struct{}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arc := strings.TrimPrefix(r.URL.Path, "/")
	if arc == "" {
		arc = "intro"
	}

	err := tmpl.Execute(w, cyoa.NextArc(arc))
	if err != nil {
		exit(err.Error())
	}
}

func main() {
	jsonFile := flag.String("json", "story.json", "json file containing the story arcs")
	usageFunc := func() {
		fmt.Fprintf(os.Stderr, "Usage: choose-your-own-adventure -json FILE\n")
		f := flag.Lookup("json")
		fmt.Fprintf(os.Stderr, "  -%s\t%s (default \"%s\")\n", f.Name, f.Usage, f.DefValue)
	}

	flag.Usage = usageFunc
	flag.Parse()

	json := readFile(*jsonFile)

	if err := cyoa.ParseJSON(json); err != nil {
		exit(fmt.Sprintf("failed to parse story file: %s", *jsonFile))
	}

	fmt.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe("localhost:8080", server{}))
}

func readFile(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		exit(fmt.Sprintf("failed to open file: %s", filename))
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		exit(fmt.Sprintf("failed to read file: %s", filename))
	}
	return bytes
}

func exit(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
