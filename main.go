package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	. "github.com/satryarangga/4venuee-api/config"
	. "github.com/satryarangga/4venuee-api/dao"
	. "github.com/satryarangga/4venuee-api/models"
)

var config = Config{}
var dao = VisitsDAO{}

// POST a new movie
func CreateVisitEndpoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var visit Visit
	if err := json.NewDecoder(r.Body).Decode(&visit); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	visit.ID = bson.NewObjectId()
	if err := dao.Insert(visit); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, visit)
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
	r.HandleFunc("/visits", CreateVisitEndpoint).Methods("POST")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}