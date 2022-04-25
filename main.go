package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var (
	DIR  string
	PORT uint
)

func init() {
	flag.StringVar(&DIR, "dir", ".", "directory to serve")
	flag.UintVar(&PORT, "port", 8090, "port to serve on")
	flag.Parse()
}

func main() {
	port := fmt.Sprintf(":%d", PORT)
	fmt.Println("http://127.0.0.1" + port)
	handlerFunc := http.HandlerFunc(handler)
	log.Fatalln(http.ListenAndServe(port, handlerFunc))
}

func handler(w http.ResponseWriter, r *http.Request) {
	rawQuery := fmt.Sprintf("?%s", r.URL.RawQuery)
	urlPath := strings.Split(r.URL.Path, "?")[0]
	filename := fmt.Sprintf("%s%s", urlPath, url.QueryEscape(rawQuery))
	paths := []string{
		filename,
		r.URL.EscapedPath(),
		r.URL.Path,
	}
	for _, path := range paths {
		name := filepath.Join(DIR, path)
		if body, err := os.ReadFile(name); err == nil {
			log.Println(name)
			writeData(w, body)
			return
		}
	}
	pattern := fmt.Sprintf("%s*", filepath.Join(DIR, urlPath))
	matches, _ := filepath.Glob(pattern)
	for _, match := range matches {
		if body, err := os.ReadFile(match); err == nil {
			log.Println(match)
			writeData(w, body)
			return
		}
	}
	log.Println(pattern)
	writeData(w, []byte("{}"))
}

func writeData(w http.ResponseWriter, data []byte) {
	contentType := "application/json; charset=utf-8"
	if !json.Valid(data) {
		contentType = http.DetectContentType(data)
	}
	w.Header().Set("Content-Type", contentType)
	_, _ = w.Write(data)
}
