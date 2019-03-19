package models

import "gopkg.in/mgo.v2/bson"

// Represents a communication, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document
type Communication struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	FirstName string        `bson:"firstName" json:"firstName"`
	LastName  string        `bson:"lastName" json:"lastName"`
	Message   string        `bson:"message" json:"message"`
	Email     string        `bson:"email" json:"email"`
	Service   string        `bson:"service" json:"service"`
	Status    string        `bson:"service" json:"service"`
}
