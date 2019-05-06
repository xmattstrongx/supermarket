package api

import (
	"encoding/json"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/xmattstrongx/supermarket/models"
)

const (
	SORTED_BY = "sortBy"
	ORDER     = "order"
	LIMIT     = "limit"
	OFFSET    = "offset"
)

type queryParameters struct {
	sortBy string
	order  string
	limit  string
	offset string
}

// ListProduce is an API handlerFunc for listing all produce inventory in the DB
func (s *Server) ListProduce(w http.ResponseWriter, r *http.Request) {
	var produce []models.Produce
	var err error
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		produce, err = s.listProduce(getQueryParams(r))
	}()
	wg.Wait()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(produce)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return
}

func (s *Server) listProduce(queryParams queryParameters) ([]models.Produce, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var produce []models.Produce
	for _, val := range s.data {
		produce = append(produce, val)
	}

	if queryParams.sortBy != "" {
		produce = sortProduce(produce, queryParams)
	}

	return produce, nil
}

func sortProduce(produce []models.Produce, queryParams queryParameters) []models.Produce {
	sortedProduce := make([]models.Produce, len(produce))

	for i, val := range produce {
		sortedProduce[i] = val
	}

	switch queryParams.sortBy {
	case "name":
		{
			if queryParams.order == "desc" {
				sort.Slice(sortedProduce, func(i, j int) bool { return sortedProduce[i].Name > sortedProduce[j].Name })
			} else {
				sort.Slice(sortedProduce, func(i, j int) bool { return sortedProduce[i].Name < sortedProduce[j].Name })
			}
		}
	case "produceCode":
		{
			if queryParams.order == "desc" {
				sort.Slice(sortedProduce, func(i, j int) bool { return sortedProduce[i].ProduceCode > sortedProduce[j].ProduceCode })
			} else {
				sort.Slice(sortedProduce, func(i, j int) bool { return sortedProduce[i].ProduceCode < sortedProduce[j].ProduceCode })
			}
		}
	case "unitPrice":
		{
			if queryParams.order == "desc" {
				sort.Slice(sortedProduce, func(i, j int) bool { return sortedProduce[i].UnitPrice > sortedProduce[j].UnitPrice })
			} else {
				sort.Slice(sortedProduce, func(i, j int) bool { return sortedProduce[i].UnitPrice < sortedProduce[j].UnitPrice })
			}
		}
	}

	return sortedProduce
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

	js, err := json.Marshal(&models.CreateProduceResponse{
		Created: validProduce,
		Invalid: invalidProduce,
	})
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

func getQueryParams(r *http.Request) queryParameters {
	query := r.URL.Query()
	return queryParameters{
		sortBy: query.Get(SORTED_BY),
		order:  query.Get(ORDER),
		limit:  query.Get(LIMIT),
		offset: query.Get(OFFSET),
	}
}
