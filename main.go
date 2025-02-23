package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"io/ioutil"
	"sync"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"gopkg.in/yaml.v2"
)

var config = struct {
	URLs map[string]string `yaml:"urls"`
}{}

var debug = os.Getenv("DEBUG")
var port = os.Getenv("PORT")
var configFile = os.Getenv("CONFIG_FILE")
var reloadInterval = os.Getenv("RELOAD_INTERVAL")
var reloadIntervalDuration time.Duration
var configMutex sync.Mutex

func init() {
	if port == "" {
		port = "8080"
	}
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

func loadConfig() error {
	configMutex.Lock()
	defer configMutex.Unlock()

	// Clear the existing configuration
	config.URLs = make(map[string]string)

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}

	return nil
}

func startConfigReloader() {
	ticker := time.NewTicker(reloadIntervalDuration)
	go func() {
		for range ticker.C {
			if err := loadConfig(); err != nil {
				log.Printf("Failed to reload config: %v", err)
			} else if debug == "true" {
				log.Printf("Config reloaded: %+v", config)
			}
		}
	}()
}

func main() {
	// Load initial configuration
	err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Start the configuration reloader
	startConfigReloader()

	// Set up the router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Define the redirect handler
	r.Get("/{link}", func(w http.ResponseWriter, r *http.Request) {
		link := chi.URLParam(r, "link")
		configMutex.Lock()
		redirect, found := config.URLs[link]
		configMutex.Unlock()
		if found {
			http.Redirect(w, r, redirect, http.StatusFound)
			return
		}
		http.NotFound(w, r)
	})

	// Start the server
	log.Printf("Starting server on :%s", port)
	log.Printf("Reload interval: %s", reloadIntervalDuration)
	http.ListenAndServe(":"+port, r)
}
