package router

import (
	"net/http"
	"sync"

	handler "github.com/sunao-uehara/hashencpassword/handlers"
)

// NewRouter registers and associates the endpoint with Handler function
// returns registered handlers
func NewRouter(wg *sync.WaitGroup) http.Handler {
	h := handler.NewHandler(wg)

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.StatsMiddleware(h.HelloHandler))
	mux.HandleFunc("/hash", h.StatsMiddleware(h.HashPostHandler))
	mux.HandleFunc("/hash/", h.StatsMiddleware(h.HashGetHandler))
	mux.HandleFunc("/stats", h.StatsGetHandler)
	mux.HandleFunc("/shutdown", h.ShutdownHandler)

	return mux
}
