package controllers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/xmattstrongx/supermarket/models"
)

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

var CreateProduce = func(w http.ResponseWriter, r *http.Request) {

	newProduce := &models.Produce{}
	err := json.NewDecoder(r.Body).Decode(newProduce) //decode the request body into struct and failed if any error occur
	if err != nil {
		Respond(w, Message(false, "Invalid request"))
		return
	}
	log.Printf("newProduce request %+v", newProduce)

	Respond(w, Message(true, "Valid request"))

	// resp := account.Create() //Create account
	// Respond(w, resp)
}
