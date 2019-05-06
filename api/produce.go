package api

import (
	"encoding/json"
	"math"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/xmattstrongx/supermarket/models"
)

// ListProduce is an API handlerFunc for listing all produce inventory in the DB
func (s *Server) ListProduce(w http.ResponseWriter, r *http.Request) {
	var b []byte
	var err error
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		b, err = s.listProduce()
	}()
	wg.Wait()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return
}

func (s *Server) listProduce() ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	b, err := json.Marshal(s.data)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// CreateProduce is an API handlerFunc for adding new produce(s) to the DB
func (s *Server) CreateProduce(w http.ResponseWriter, r *http.Request) {
	newProduce := &[]models.Produce{}
	err := json.NewDecoder(r.Body).Decode(newProduce) //decode the request body into struct and failed if any error occur
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validProduce := []models.Produce{}
	invalidProduce := []models.Produce{}
	wg := sync.WaitGroup{}

	for _, val := range *newProduce {
		if _, exists := s.data[strings.ToUpper(val.ProduceCode)]; exists {
			invalidProduce = append(invalidProduce, val)
			break
		}

		if !isValidProduceCode(val.ProduceCode) {
			invalidProduce = append(invalidProduce, val)
			break
		}

		wg.Add(1)
		go func(val models.Produce) {
			defer wg.Done()
			p, err := s.createProduce(val)
			if err != nil {
				invalidProduce = append(invalidProduce, val)
				return
			}
			validProduce = append(validProduce, p)
		}(val)
	}
	wg.Wait()

	resp := &models.CreateProduceReponse{
		Created: validProduce,
		Invalid: invalidProduce,
	}

	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status := http.StatusCreated
	switch {
	case len(invalidProduce) > 0 && len(validProduce) > 0:
		status = http.StatusMultiStatus
	case len(invalidProduce) > 0 && len(validProduce) == 0:
		// TODO decide if an error should be returned for each specific failure
		status = http.StatusBadRequest
	default:
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return
}

// createProduce cannot return an error in this implementation
// however this method is part of an interface so a future implementation that uses
// a real database would need the option to return an error
func (s *Server) createProduce(newProduce models.Produce) (models.Produce, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[strings.ToUpper(newProduce.ProduceCode)] = models.Produce{
		Name:        newProduce.Name,
		ProduceCode: strings.ToUpper(newProduce.ProduceCode),
		UnitPrice:   math.Round(newProduce.UnitPrice*100) / 100,
	}
	return newProduce, nil
}

// DeleteProduce is an API handlerFunc for adding removing produce from the DB
func (s *Server) DeleteProduce(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := s.deleteProduce(vars["productCode"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}

// deleteProduce cannot return an error in this implementation
// however this method is part of an interface so a future implementation that uses
// a real database would need the option to return an error
func (s *Server) deleteProduce(produceCode string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.data, produceCode)
	return nil
}

func isValidProduceCode(produceCode string) bool {
	validProduceCode := regexp.MustCompile(`[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{4}`)
	return validProduceCode.MatchString(produceCode)
}
