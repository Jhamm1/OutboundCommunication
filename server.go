package main

import (
	// Standard library packages

	"fmt"
	"net/http"

	mgo "gopkg.in/mgo.v2"

	// Third party packages
	"github.com/julienschmidt/httprouter"
	"github.com/swhite24/go-rest-tutorial/controllers"
)

// var (
// 	uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
// 	exchangeName = flag.String("exchange", "test-exchange", "Durable AMQP exchange name")
// 	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
// 	routingKey   = flag.String("key", "test-key", "AMQP routing key")
// 	body         = flag.String("body", "foobar", "Body of message")
// 	reliable     = flag.Bool("reliable", true, "Wait for the publisher confirmation before exiting")
// )

// func init() {
// 	flag.Parse()
// }

func main() {
	// Instantiate a new router
	r := httprouter.New()

	// Get a UserController instance
	uc := controllers.NewUserController(getSession())

	// Get a user resource
	r.GET("/user/:id", uc.GetUser)

	r.POST("/user", uc.CreateUser)

	r.DELETE("/user/:id", uc.RemoveUser)

	// Fire up the server
	http.ListenAndServe("localhost:3001", r)
	fmt.Println("localhost:3001")

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
