package app

import (
	"Project/internal/manager"
	"Project/internal/repository"
	"fmt"
	"net/http"
)

type HezzlWebService struct {
	config  *Config
	manager *manager.Manager
	router  *http.ServeMux
	store   *repository.Store
}

func New(config *Config) *HezzlWebService {
	return &HezzlWebService{
		config: config,
	}
}

func (s *HezzlWebService) Start() error {
	if err := s.configureStore(); err != nil {
		fmt.Println(err)
		return err
	}
	s.configureRoutes()
	s.configureManager()

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *HezzlWebService) configureStore() error {
	if s.store == nil {
		st := &repository.Store{}
		st.ConfigureStore()

		s.store = st
	}
	return nil
}

func (s *HezzlWebService) configureRoutes() {
	if s.router == nil {
		s.router = http.NewServeMux()
		s.ConfigureRouter()
	}
}

func (s *HezzlWebService) configureManager() {
	if s.manager == nil {
		s.manager = manager.NewManager(s.store)
	}
}
