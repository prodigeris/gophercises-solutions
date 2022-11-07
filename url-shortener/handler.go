package urlshort

import (
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	mux := http.NewServeMux()
	for path, url := range pathsToUrls {
		mux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
			http.Redirect(writer, request, url, 308)
		})
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		h, s := mux.Handler(request)
		if s != "" {
			h.ServeHTTP(writer, request)
		} else {
			fallback.ServeHTTP(writer, request)
		}
	}
}

type Route struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
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
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	routes := make([]Route, 0)
	if err := yaml.Unmarshal(yml, &routes); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	for _, route := range routes {
		mux.HandleFunc(route.Path, func(writer http.ResponseWriter, request *http.Request) {
			http.Redirect(writer, request, route.Url, 308)
		})
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		h, s := mux.Handler(request)
		if s != "" {
			h.ServeHTTP(writer, request)
		} else {
			fallback.ServeHTTP(writer, request)
		}
	}, nil
}
