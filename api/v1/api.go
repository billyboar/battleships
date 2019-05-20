package v1

import (
	"github.com/billyboar/battleships/models/db"
	"github.com/gorilla/mux"
)

// APIServer contains db connection and api routes
type APIServer struct {
	Router *mux.Router
	Store  *db.Store
}

// NewAPIServer creates new server struct with redis connection
func NewAPIServer() (*APIServer, error) {
	store, err := db.NewStore()
	if err != nil {
		return nil, err
	}

	return &APIServer{
		Router: mux.NewRouter(),
		Store:  store,
	}, nil
}

// RegisterRoutes adds new routes to main routes handler
func (s *APIServer) RegisterRoutes() {
	s.Router.Use(s.GlobalCORSMiddleware)

	s.Router.HandleFunc("/health", HealthCheck).Methods("GET")

	apiRoute := s.Router.PathPrefix("/api/v1").Subrouter()

	s.LoadSessionRoutes(apiRoute)
}
