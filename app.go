package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/Jhamm1/OutboundCommunication/controllers"
	. "github.com/Jhamm1/OutboundCommunication/testpage"
	"github.com/gorilla/mux"
)

func apiEndpoints() {
	r := mux.NewRouter()
	r.HandleFunc("/communications", controllers.AllCommunicationsEndPoint).Methods("GET")
	r.HandleFunc("/communications", controllers.CreateCommunicationEndPoint).Methods("POST")
	r.HandleFunc("/communications", controllers.UpdateCommunicationEndPoint).Methods("PUT")
	r.HandleFunc("/communications", controllers.DeleteCommunicationEndPoint).Methods("DELETE")
	r.HandleFunc("/communications/{id}", controllers.FindCommunicationEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3002", r); err != nil {
		log.Fatal(err)
	}
}

// HTTP request routes
func main() {
	apiEndpoints()
	if runtime.GOOS == "windows" {
		fmt.Println("Can't Execute this on a windows machine")
	} else {
		testpage.execute()

	}
	fmt.Println("Can't Execute this on a windows machine")
	testpage.openbrowser("")
}
