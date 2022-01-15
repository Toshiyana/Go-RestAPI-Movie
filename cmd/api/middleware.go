package main

import "net/http"

func (app *application) enableCORS(next http.Handler) http.Handler {
	// CORS: Corss-Origin Resource Sharing
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// *: Allow all requests
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}
