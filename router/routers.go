package router

import (
	"bullshape/ctrls"

	"github.com/gorilla/mux"
)

func apiRoutes(router *mux.Router) {

	// usersRoutes := router.PathPrefix("/api/user").Subrouter()
	// // usersRoutes := mux.NewRouter().PathPrefix("/api/users").Subrouter()
	// usersRoutes.HandleFunc("", ctrl.GetUser).Methods("GET")
	// usersRoutes.HandleFunc("", ctrl.CreateUser).Methods("POST")
	// usersRoutes.HandleFunc("", ctrl.DeleteUser).Methods("DELETE")
	router.HandleFunc("/api/companies", ctrls.CreateCompany).Methods("POST")

	companiesRoutes := router.PathPrefix("/api/companies/{id}").Subrouter()
	companiesRoutes.HandleFunc("", ctrls.GetCompany).Methods("GET")
	companiesRoutes.HandleFunc("", ctrls.DeleteCompany).Methods("DELETE")
	companiesRoutes.HandleFunc("", ctrls.UpdateCompany).Methods("PATCH")

}

// func authRoutes(router *mux.Router) {
//
// }

func LoadRoutes() *mux.Router {
	routes := mux.NewRouter()
	// routes.HandleFunc("/aliveness", ctrl.Aliveness)
	apiRoutes(routes)
	// authRoutes(routes)
	return routes
}
