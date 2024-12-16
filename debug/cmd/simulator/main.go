package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type RegionHealth struct {
	mu     sync.RWMutex
	health map[string]bool
}

func main() {
	regions := map[string]string{
		"west-us":    "http://west-us:8080",
		"central-us": "http://central-us:8080",
	}

	health := &RegionHealth{
		health: make(map[string]bool),
	}

	// Health check routine
	// Health check routine
	go func() {
		for {
			for region, url := range regions {
				resp, err := http.Get(url + "/api/health/live") // Changed from /health to /api/health/live
				health.mu.Lock()
				health.health[region] = err == nil && resp != nil && resp.StatusCode == 200
				health.mu.Unlock()
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// Proxy handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		health.mu.RLock()
		primaryHealthy := health.health["west-us"]
		health.mu.RUnlock()

		target := regions["west-us"]
		if !primaryHealthy {
			target = regions["central-us"]
		}

		targetURL, _ := url.Parse(target)
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.ServeHTTP(w, r)
	})

	log.Printf("Traffic Manager simulator starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
