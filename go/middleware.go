package main

import (
	"log"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Printf("[web] %s: Begin of %s\n", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		log.Printf("[web] %s: %s\n", r.Method, r.RequestURI)
	})
}
func cachingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Cache Handler started")
		onlineCache(r.RequestURI)
		log.Println("Cache Handler finished")
		next.ServeHTTP(w, r)
	})
}


