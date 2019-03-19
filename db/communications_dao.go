package dao

import (
	"log"

	. "github.com/Jhamm1/OutboundCommunication/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CommunicationsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "communications"
)

// Establish a connection to database
func (m *CommunicationsDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of communications
func (m *CommunicationsDAO) FindAll() ([]Communication, error) {
	var communications []Communication
	err := db.C(COLLECTION).Find(bson.M{}).All(&communications)
	return communications, err
}

// Find a communication by its id
func (m *CommunicationsDAO) FindById(id string) (Communication, error) {
	var communication Communication
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&communication)
	return communication, err
}

// Insert a communication into database
func (m *CommunicationsDAO) Insert(communication Communication) error {
	err := db.C(COLLECTION).Insert(&communication)
	return err
}

// Delete an existing communication
func (m *CommunicationsDAO) Delete(communication Communication) error {
	err := db.C(COLLECTION).Remove(&communication)
	return err
}

// Update an existing communication
func (m *CommunicationsDAO) Update(communication Communication) error {
	err := db.C(COLLECTION).UpdateId(communication.ID, &communication)
	return err
}

// Find a communication by it
func (m *CommunicationsDAO) FindandUpdate(selector interface{}, update interface{}) (Communication, error) {
	var communication Communication
	err := db.C(COLLECTION).Update(selector, update)
	// err2 := db.C(COLLECTION).UpdateId(bson.ObjectIdHex(id), &communication)
	return communication, err
}
