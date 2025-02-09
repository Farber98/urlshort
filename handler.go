package urlshort

import (
	"encoding/json"
	"log"
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
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if url, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pu := parseYaml(yml)
	if len(pu) > 0 {
		puMap := arrToMap(pu)
		return MapHandler(puMap, fallback), nil
	}
	return MapHandler(nil, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	urls := parseJson(json)
	if len(urls.Urls) > 0 {
		puMap := arrToMap(urls.Urls)
		return MapHandler(puMap, fallback), nil
	}
	return YAMLHandler(nil, fallback)
}

func parseYaml(yamlFile []byte) (pu []pathUrls) {
	err := yaml.Unmarshal(yamlFile, &pu)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return
}

func arrToMap(puArr []pathUrls) map[string]string {
	mp := make(map[string]string)
	for _, elem := range puArr {
		mp[elem.Path] = elem.Url
	}
	return mp
}

type pathUrls struct {
	Path string `yaml:"path,omitempty" json:"path,omitempty"`
	Url  string `yaml:"url,omitempty" json:"url,omitempty"`
}

func parseJson(jsonFile []byte) (urls urls) {
	err := json.Unmarshal(jsonFile, &urls)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return
}

type urls struct {
	Urls []pathUrls `json:"urls,omitempty"`
}
