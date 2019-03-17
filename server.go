package main

import (
	// Standard library packages

	"net/http"

	mgo "gopkg.in/mgo.v2"

	// Third party packages
	"github.com/Jhamm1/OutboundCommunication/controllers"
	"github.com/julienschmidt/httprouter"
	//"github.com/swhite24/go-rest-tutorial/controllers"
)

func main() {

	// Instantiate a new router
	r := httprouter.New()

	// Get a UserController instance
	uc := controllers.NewCommunicationController(getSession())

	// Get a user resource
	r.GET("/user/:id", uc.GetCommunication)

	// Create a new user
	r.POST("/user", uc.CreateCommunication)

	// Remove an existing user
	r.DELETE("/user/:id", uc.RemoveCommunication)

	// Fire up the server
	http.ListenAndServe("localhost:3001", r)

}

// getSession creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}

	// Deliver session
	return s
}