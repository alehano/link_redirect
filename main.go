package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/configor"
)

// Config structure to hold the URLs
type Config struct {
	URLs []struct {
		URL      string `yaml:"url"`
		Redirect string `yaml:"redirect"`
	} `yaml:"urls"`
}

var config Config

func main() {
	// Load configuration with auto-reload enabled
	err := configor.New(&configor.Config{
		AutoReload:         true,
		AutoReloadInterval: 10 * time.Second,
	}).Load(&config, "config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set up the router
	r := chi.NewRouter()

	// Define the redirect handler
	r.Get("/{link}", func(w http.ResponseWriter, r *http.Request) {
		link := chi.URLParam(r, "link")
		for _, u := range config.URLs {
			if u.URL == link {
				http.Redirect(w, r, u.Redirect, http.StatusFound)
				return
			}
		}
		http.NotFound(w, r)
	})

	// Start the server
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
