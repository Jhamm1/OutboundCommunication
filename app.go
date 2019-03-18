package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	. "github.com/mlabouardy/movies-restapi/config"
	. "github.com/mlabouardy/movies-restapi/dao"
	. "github.com/mlabouardy/movies-restapi/models"
	"github.com/streadway/amqp"
)

var config = Config{}
var dao = CommunicationsDAO{}

// GET list of communications
func AllMoviesEndPoint(w http.ResponseWriter, r *http.Request) {
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

// POST a new communication
func CreateCommunicationEndPoint(w http.ResponseWriter, r *http.Request) {
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

	///------------------------ Rabbit MQ -----------------------------------//
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Outbound_communication_service", // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := communication
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			//Body:        []byte(body),
			Body: []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")

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

// Define HTTP request routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/communications", AllMoviesEndPoint).Methods("GET")
	r.HandleFunc("/communications", CreateMovieEndPoint).Methods("POST")
	r.HandleFunc("/communications", UpdateMovieEndPoint).Methods("PUT")
	r.HandleFunc("/communications", DeleteMovieEndPoint).Methods("DELETE")
	r.HandleFunc("/communications/{id}", FindMovieEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3002", r); err != nil {
		log.Fatal(err)
	}
}
