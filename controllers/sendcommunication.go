package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Jhamm1/OutboundCommunication/models"
	"github.com/julienschmidt/httprouter"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"github.com/swhite24/go-rest-tutorial/models"
)

// type (
// 	// CommunicationController represents the controller for operating on the Communication resource
// 	CommunicationController struct{}
// )

type (
	// CommunicationController represents the controller for operating on the CommunicationController resource
	CommunicationController struct {
		session *mgo.Session
	}
)

// NewCommunicationController provides a reference to a CommunicationController with provided mongo session
func NewCommunicationController(s *mgo.Session) *CommunicationController {
	return &CommunicationController{}
}

// GetCommunication retrieves an individual user resource
func (uc CommunicationController) GetCommunication(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an example user
	// u := models.User{
	// 	FirstName: "John",
	// 	LastName:  "Smith",
	// 	Message:   "This is a test Bob Smith",
	// 	Email:     "John.Smith@email.com",
	// 	Service:   "Security advice",
	// }

	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Stub user
	u := models.Communication{}

	// Fetch user
	if err := uc.session.DB("go_rest_tutorial").C("users").FindId(oid).One(&u); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// CreateCommunication creates a new communication resource
func (uc CommunicationController) CreateCommunication(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an user to be populated from the body
	u := models.Communication{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.Id = bson.NewObjectId()

	// Write the user to mongo
	uc.session.DB("Outbound_communication").C("communication").Insert(u)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)

}

// RemoveCommunication removes an existing communication resource
func (uc CommunicationController) RemoveCommunication(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove user
	if err := uc.session.DB("Outbound_communication").C("communication").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}
