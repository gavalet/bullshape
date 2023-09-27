package router

import (
	"bullshape/ctrls"

	"github.com/gorilla/mux"
)

func (s *Server) apiRoutes(router *mux.Router) {

	ctrl := ctrls.NewCtrlServices(s.Logger, s.DB, s.Kafka)

	router.HandleFunc("/api/companies", ctrl.CreateCompany).Methods("POST")
	companiesRoutes := router.PathPrefix("/api/companies/{id}").Subrouter()
	companiesRoutes.HandleFunc("", ctrl.GetCompany).Methods("GET")
	companiesRoutes.HandleFunc("", ctrl.DeleteCompany).Methods("DELETE")
	companiesRoutes.HandleFunc("", ctrl.UpdateCompany).Methods("PATCH")

}

func (s *Server) authRoutes(router *mux.Router) {
	ctrl := ctrls.NewCtrlServices(s.Logger, s.DB, s.Kafka)
	router.HandleFunc("/api/user", ctrl.CreateUser).Methods("POST")
	router.HandleFunc("/api/user/login", ctrl.Authenticate).Methods("POST")
}

func (s *Server) loadRoutes() *mux.Router {
	routes := mux.NewRouter()
	// routes.HandleFunc("/aliveness", ctrl.Aliveness)
	s.apiRoutes(routes)
	s.authRoutes(routes)
	routes.Use(s.jwtAuthentication)
	routes.Use(s.recovery)

	return routes
}
