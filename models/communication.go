package models

import "gopkg.in/mgo.v2/bson"

// Represents a movie, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document
// type Movie struct {
// 	ID          bson.ObjectId `bson:"_id" json:"id"`
// 	Name        string        `bson:"name" json:"name"`
// 	CoverImage  string        `bson:"cover_image" json:"cover_image"`
// 	Description string        `bson:"description" json:"description"`
// }

type Communication struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	FirstName string        `bson:"firstName" json:"firstName"`
	LastName  string        `bson:"lastName" json:"lastName"`
	Message   string        `bson:"message" json:"message"`
	Email     string        `bson:"email" json:"email"`
	Service   string        `bson:"service" json:"service"`
}
