package api

import (
	"context"
	"net"
	"net/http"

	"github.com/aler9/rtsp-simple-server/internal/logger"
)

func muxHandle(mux *http.ServeMux, method string, path string,
	cb func(w http.ResponseWriter, r *http.Request)) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		cb(w, r)
	})
}

// Parent is implemented by program.
type Parent interface {
	Log(logger.Level, string, ...interface{})
}

// API is an HTTP API server.
type API struct {
	s *http.Server
}

// New allocates an API.
func New(
	address string,
	parent Parent,
) (*API, error) {
	a := &API{}

	ln, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	muxHandle(mux, "GET", "/config", a.onGETConfig)

	// POST /config
	// GET /config/paths
	// GET /config/path/id
	// POST /config/path/id
	// GET /paths
	// GET /clients
	// POST /client/id/kick

	a.s = &http.Server{
		Handler: mux,
	}

	go a.s.Serve(ln)

	parent.Log(logger.Info, "[api] listener opened on "+address)

	return a, nil
}

// Close closes all the API resources.
func (a *API) Close() {
	a.s.Shutdown(context.Background())
}

func (a *API) onGETConfig(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("TODO\n"))
}
