package router

import (
	"github.com/miledxz/salesforge/handler"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/sequences", handler.CreateSequence).Methods("POST")
	r.HandleFunc("/steps/{id}", handler.UpdateStep).Methods("PUT")
	r.HandleFunc("/steps/{id}", handler.DeleteStep).Methods("DELETE")
	r.HandleFunc("/sequences/{id}/tracking", handler.UpdateTracking).Methods("PUT")
	return r
}
