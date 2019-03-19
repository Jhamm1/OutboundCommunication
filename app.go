package main

import (
	"log"
	"net/http"

	"github.com/Jhamm1/OutboundCommunication/controllers"
	"github.com/darkowlzz/openurl"
	"github.com/gorilla/mux"
)

// HTTP request routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/communications", controllers.AllCommunicationsEndPoint).Methods("GET")
	r.HandleFunc("/communications", controllers.CreateCommunicationEndPoint).Methods("POST")
	r.HandleFunc("/communications", controllers.UpdateCommunicationEndPoint).Methods("PUT")
	r.HandleFunc("/communications", controllers.DeleteCommunicationEndPoint).Methods("DELETE")
	r.HandleFunc("/communications/{id}", controllers.FindCommunicationEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3002", r); err != nil {
		log.Fatal(err)
	}

	if err := openurl.Open("http://example.com"); err != nil {
		log.Fatal(err)
	}
}
