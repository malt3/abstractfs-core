package http

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/malt3/abstractfs-core/api"
	"github.com/malt3/abstractfs-core/sri"
)

// Handler is a CAS http handler.
// It implements the CAS http protocol.
// It forwards requests to a CAS backend.
type Handler struct {
	cas api.CAS
}

func NewHandler(cas api.CAS) http.Handler {
	return &Handler{
		cas: cas,
	}
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		s.handleGet(w, req)
	case http.MethodPut:
		s.handlePut(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGet handles a GET request.
// It expects the sri in the following format:
// /cas/<hash-function>/<hash-value-hex>
func (s *Handler) handleGet(w http.ResponseWriter, req *http.Request) {
	integrity, err := parsePath(req.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err := s.cas.Open(integrity.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer body.Close()
	if _, err := io.Copy(w, body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Handler) handlePut(w http.ResponseWriter, req *http.Request) {
	integrity, err := parsePath(req.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.cas.Write(integrity.String(), req.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// parsePath parses the path and returns the sri.
// It expects the sri in the following format:
// /cas/<hash-function>/<hash-value-hex>
func parsePath(path string) (sri.Integrity, error) {
	if !strings.HasPrefix(path, "/cas/") {
		return sri.Integrity{}, errors.New("invalid path: must start with /cas/")
	}
	path = path[len("/cas/"):]
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		return sri.Integrity{}, errors.New("invalid path: must have format /cas/<hash-function>/<hash-value-hex>")
	}
	alg, err := sri.AlgorithmFromString(parts[0])
	if err != nil {
		return sri.Integrity{}, fmt.Errorf("invalid path: %w", err)
	}
	hash, err := hex.DecodeString(parts[1])
	if err != nil {
		return sri.Integrity{}, fmt.Errorf("invalid path: %w", err)
	}
	return sri.Integrity{
		Algorithm: alg,
		Hash:      hash,
	}, nil
}
