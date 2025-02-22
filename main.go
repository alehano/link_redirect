package main

import (
	"log"
	"net/http"
	"os"
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
var configFile = os.Getenv("CONFIG_FILE")
var reloadInterval = os.Getenv("RELOAD_INTERVAL")
var reloadIntervalDuration time.Duration

func init() {
	if configFile == "" {
		configFile = "config.yml"
	}
	if reloadInterval == "" {
		reloadInterval = "10s"
	}
	ri, err := time.ParseDuration(reloadInterval)
	if err != nil {
		log.Fatalf("Failed to parse RELOAD_INTERVAL: %v", err)
	}
	reloadIntervalDuration = ri
}

func main() {
	// Load configuration with auto-reload enabled
	err := configor.New(&configor.Config{
		AutoReload:         true,
		AutoReloadInterval: reloadIntervalDuration,
	}).Load(&config, configFile)
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
