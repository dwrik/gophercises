package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		url, ok := pathsToUrls[r.URL.Path]
		if !ok {
			fallback.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
	return handler
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml := []mapping{}
	err := yaml.Unmarshal(yml, &parsedYaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
// [
//   {
//     "path": "programming",
//     "url": "https://www.reddit.com/r/learnprogramming"
//   },
//   ...
// ]
//
// The only errors that can be returned all related to having
// invalid JSON data.

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON := []mapping{}
	err := json.Unmarshal(jsn, &parsedJSON)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJSON)
	return MapHandler(pathMap, fallback), nil
}

// Helper method
// for building map

func buildMap(parsedData []mapping) map[string]string {
	pathMap := map[string]string{}
	for _, m := range parsedData {
		pathMap[m.Path] = m.Url
	}
	return pathMap
}

// Custom type for
// mapping YAML and JSON

type mapping struct {
	Path string
	Url  string
}
