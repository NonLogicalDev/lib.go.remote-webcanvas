package wscanvas

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"

	"github.com/NonLogicalDev/lib.go.remote-webcanvas/wscanvas/static"
)

type staticHandler struct {
	webSocketAddress string
}

func (s *staticHandler) expandVariables(data []byte) []byte {
	vars := map[string]string{
		"WEB_SOCKET_ADDRESS": s.webSocketAddress,
	}
	for name, value := range vars {
		data = bytes.ReplaceAll(data, []byte("%%" + name + "%%"), []byte(value))
	}
	return data
}

func (s *staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v", r.URL)

	slug := chi.URLParam(r, "*")
	if slug == "" {
		slug = "index.html"
	}
	ext := filepath.Ext(slug)

	data, err := static.Asset(slug)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error %v", err)
	}

	w.Header().Set("Content-Type", "text/html")
	if ext == "js" {
		w.Header().Set("Content-Type", "text/javascript")
	}

	w.Write(s.expandVariables(data))
}
