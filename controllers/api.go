package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/mgo.v2/bson"

	. "github.com/Jhamm1/OutboundCommunication/config"
	. "github.com/Jhamm1/OutboundCommunication/db"
	. "github.com/Jhamm1/OutboundCommunication/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var config = Config{}
var dao = CommunicationsDAO{}

// GET list of communications
func AllCommunicationsEndPoint(w http.ResponseWriter, r *http.Request) {
	communications, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, communications)
}

// GET a communication by its ID
func FindCommunicationEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	communication, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Communication ID")
		return
	}
	respondWithJson(w, http.StatusOK, communication)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Function for publishing msg onto the queue
func publishCommunicationOnQueue(connection string, queueName string, msg Communication) {

	conn, err := amqp.Dial(connection)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	e, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	body := e
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")

}

// POST a new communication
func CreateCommunicationEndPoint(w http.ResponseWriter, r *http.Request) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var queueURL = os.Getenv("QUEUEURL")
	var queueName = os.Getenv("QUEUENAME")
	log.Print("Test: ", queueURL)
	// var yourDomain string = os.Getenv("DOMAIN") // e.g. mg.yourcompany.com

	// // You can find the Private API Key in your Account Menu, under "Settings":
	// // (https://app.mailgun.com/app/account/security)
	// var privateAPIKey string = os.Getenv("APIKEY")

	defer r.Body.Close()
	var communication Communication
	if err := json.NewDecoder(r.Body).Decode(&communication); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	communication.ID = bson.NewObjectId()
	if err := dao.Insert(communication); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, communication)

	//ToDo
	//Update status in DB to 'processingRequest'
	//user := &Communication{"": }
	b, err := json.Marshal(communication.Email)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	///------------------------ Rabbit MQ -----------------------------------//
	publishCommunicationOnQueue(queueURL, queueName, communication)

}

// PUT update an existing communication
func UpdateCommunicationEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var communication Communication
	if err := json.NewDecoder(r.Body).Decode(&communication); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(communication); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing communication
func DeleteCommunicationEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var communication Communication
	if err := json.NewDecoder(r.Body).Decode(&communication); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(communication); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}
