package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"golang.org/x/net/http2" // HTTP/2 support
)

// Middleware to log requests
func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s from %s\n", r.Method, r.URL, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// Create a reverse proxy with response modification and HTTP/2 support
func newReverseProxy(target string) *httputil.ReverseProxy {
	// Parse the target URL
	url, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Could not parse target URL: %v", err)
	}

	// Create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Modify the response
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("X-Proxy-App", "GoProxy")
		// Example: Add a footer to all HTML responses
		if resp.Header.Get("Content-Type") == "text/html" {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			modifiedBody := string(body) + "<!-- Footer added by Go Proxy -->"
			resp.Body = io.NopCloser(io.Reader(strings.NewReader(modifiedBody)))
			resp.ContentLength = int64(len(modifiedBody))
			resp.Header.Set("Content-Length", string(len(modifiedBody)))
		}
		return nil
	}
	return proxy
}

// Set up the HTTP/2 server for Go 1.22
func startHTTP2Server(proxy http.Handler) {
	server := &http.Server{
		Addr:    ":8080",
		Handler: logRequestMiddleware(proxy), // Wrap the proxy handler with logging middleware
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	// Enable HTTP/2 support
	http2Server := &http2.Server{}
	err := http2.ConfigureServer(server, http2Server)
	if err != nil {
		log.Fatalf("Error configuring HTTP/2 server: %v", err)
	}

	log.Printf("Starting HTTP/2 proxy server on https://localhost:8080")

	// Start the server (use TLS for HTTP/2)
	log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}

func main() {
	// Define the target service you want to proxy to
	target := "https://jsonplaceholder.typicode.com"

	// Create a reverse proxy instance
	proxy := newReverseProxy(target)

	// Start the HTTP/2 server
	startHTTP2Server(proxy)
}
