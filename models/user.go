package models

type (
	// User represents the structure of our resource
	User struct {
		//Id        bson.ObjectId `json:"id" bson:"_id"`
		FirstName string `json:"firstName" bson:"firstName"`
		LastName  string `json:"lastName" bson:"lastName"`
		Message   string `json:"message" bson:"message"`
		Email     string `json:"email" bson:"email"`
		Service   string `json:"service" bson:"service"`
	}
)
