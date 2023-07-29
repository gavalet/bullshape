package router

import (
	"bullshape/ctrls"

	"github.com/gorilla/mux"
)

func apiRoutes(router *mux.Router) {
	router.HandleFunc("/api/companies", ctrls.CreateCompany).Methods("POST")

	companiesRoutes := router.PathPrefix("/api/companies/{id}").Subrouter()
	companiesRoutes.HandleFunc("", ctrls.GetCompany).Methods("GET")
	companiesRoutes.HandleFunc("", ctrls.DeleteCompany).Methods("DELETE")
	companiesRoutes.HandleFunc("", ctrls.UpdateCompany).Methods("PATCH")

}

func authRoutes(router *mux.Router) {
	router.HandleFunc("/api/user", ctrls.CreateUser).Methods("POST")
	router.HandleFunc("/api/user/login", ctrls.Authenticate).Methods("POST")
}

func LoadRoutes() *mux.Router {
	routes := mux.NewRouter()
	// routes.HandleFunc("/aliveness", ctrl.Aliveness)
	apiRoutes(routes)
	authRoutes(routes)
	routes.Use(jwtAuthentication)

	return routes
}
