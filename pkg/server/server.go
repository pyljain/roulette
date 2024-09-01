package server

import (
	"fmt"
	"log"
	"net/http"
	"roulette/internal/config"
	"roulette/internal/lru"
	"roulette/internal/storage"
)

type Server struct {
	config  *config.Config
	cache   lru.LRUCache
	storage storage.Storage
}

func New(config *config.Config, cache lru.LRUCache, storage storage.Storage) *Server {
	return &Server{
		config:  config,
		cache:   cache,
		storage: storage,
	}
}

func (s *Server) Run() error {
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.config.Port), s)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	path = path[1:]
	log.Printf("PATH is %s", path)

	contents, err := s.cache.Get(path)
	if err == nil {
		fmt.Printf("Cache hit for %s\n", path)
		w.WriteHeader(http.StatusOK)
		w.Write(contents)
		return
	}

	// Check for file in bucket
	fmt.Printf("Cache miss for %s\n", path)
	contents, err = s.storage.Read(r.Context(), path)

	if err != nil {
		log.Printf("Error when reading contents from GCS %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(contents)

	// Hydrate cache async
	go func() {
		err = s.cache.Add(path, contents)
		if err != nil {
			log.Printf("Unable to add item to cache: %s", path)
			return
		}
	}()

}
