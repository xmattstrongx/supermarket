package api

import (
	"encoding/json"
	"math"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/xmattstrongx/supermarket/models"
)

func (s *Server) listProduce(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(s.data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
	return
}

func (s *Server) createProduce(w http.ResponseWriter, r *http.Request) {

	newProduce := &models.Produce{}
	err := json.NewDecoder(r.Body).Decode(newProduce) //decode the request body into struct and failed if any error occur
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// verify the user is not trying to override an already existing item
	if _, exists := s.data[strings.ToUpper(newProduce.ProduceCode)]; exists {
		http.Error(w, "Unable to add already existing item", http.StatusBadRequest)
		return
	}

	if !isValidProduceCode(newProduce.ProduceCode) {
		http.Error(w, "Invalid product code", http.StatusBadRequest)
		return
	}

	// TODO make actual insert
	s.data[strings.ToUpper(newProduce.ProduceCode)] = models.Produce{
		Name:        newProduce.Name,
		ProduceCode: strings.ToUpper(newProduce.ProduceCode),
		UnitPrice:   math.Round(newProduce.UnitPrice*100) / 100,
	}

	js, err := json.Marshal(newProduce)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
	return
}

func (s *Server) deleteProduce(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	delete(s.data, vars["productCode"])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}

func isValidProduceCode(produceCode string) bool {
	validProduceCode := regexp.MustCompile(`[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{4}`)
	return validProduceCode.MatchString(produceCode)
}
