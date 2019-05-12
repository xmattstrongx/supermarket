package api

import (
	"encoding/json"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/xmattstrongx/supermarket/models"
)

const (
	SORTED_BY = "sort_by"
	ORDER     = "order"
	LIMIT     = "limit"
	OFFSET    = "offset"

	QUERY_PARAM_NAME         = "name"
	QUERY_PARAM_PRODUCE_CODE = "producecode"
	QUERY_PARAM_UNIT_PRICE   = "unitprice"
	QUERY_PARAM_DESC         = "desc"
	QUERY_PARAM_DESCENDING   = "descending"
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

func getQueryParams(r *http.Request) queryParameters {
	query := r.URL.Query()
	return queryParameters{
		sortBy: query.Get(SORTED_BY),
		order:  query.Get(ORDER),
		limit:  query.Get(LIMIT),
		offset: query.Get(OFFSET),
	}
}

func (s *Server) listProduce(queryParams queryParameters) ([]models.Produce, error) {
	// s.mutex.RLock()
	// defer s.mutex.RUnlock()

	var produce []models.Produce
	for _, val := range s.data {
		s.mutex.Lock()
		produce = append(produce, val)
		s.mutex.Unlock()
	}

	if queryParams.sortBy != "" {
		produce = sortProduce(produce, queryParams)
	}

	return produce, nil
}

// sortProduce will attempt to sort the produce slice passed in based on the query parameters
// provided by the request.
func sortProduce(produce []models.Produce, queryParams queryParameters) []models.Produce {
	sortedProduce := make([]models.Produce, len(produce))

	for i, val := range produce {
		sortedProduce[i] = val
	}

	order := queryParams.order
	sortBy := queryParams.sortBy
	// if sortBy is not passed default to name
	if sortBy == "" {
		sortBy = QUERY_PARAM_NAME
	}

	switch strings.ToLower(queryParams.sortBy) {
	case QUERY_PARAM_NAME:
		{
			if order == QUERY_PARAM_DESC || order == QUERY_PARAM_DESCENDING {
				sort.Slice(sortedProduce, func(i, j int) bool { return sortedProduce[i].Name > sortedProduce[j].Name })
			} else {
				sort.Slice(sortedProduce, func(i, j int) bool { return sortedProduce[i].Name < sortedProduce[j].Name })
			}
		}
	case QUERY_PARAM_PRODUCE_CODE:
		{
			if order == QUERY_PARAM_DESC || order == QUERY_PARAM_DESCENDING {
				sort.Slice(sortedProduce, func(i, j int) bool { return sortedProduce[i].ProduceCode > sortedProduce[j].ProduceCode })
			} else {
				sort.Slice(sortedProduce, func(i, j int) bool { return sortedProduce[i].ProduceCode < sortedProduce[j].ProduceCode })
			}
		}
	case QUERY_PARAM_UNIT_PRICE:
		{
			if order == QUERY_PARAM_DESC || order == QUERY_PARAM_DESCENDING {
				sort.Slice(sortedProduce, func(i, j int) bool { return sortedProduce[i].UnitPrice > sortedProduce[j].UnitPrice })
			} else {
				sort.Slice(sortedProduce, func(i, j int) bool { return sortedProduce[i].UnitPrice < sortedProduce[j].UnitPrice })
			}
		}
	}

	// first split the produce by the offset
	offset, err := strconv.ParseInt(queryParams.offset, 10, 0)
	if err != nil || offset >= int64(len(sortedProduce)) || offset < 0 {
		offset = 0
	}

	sortedProduce = sortedProduce[offset:]

	// next split the produce so that only the limit requested is returned
	limit, err := strconv.ParseInt(queryParams.limit, 10, 0)
	if err != nil || limit > int64(len(sortedProduce)) || limit < 0 {
		limit = int64(len(sortedProduce))
	}

	return sortedProduce[:limit]
}

// CreateProduce is an API handlerFunc for adding new produce(s) to the DB
func (s *Server) CreateProduce(w http.ResponseWriter, r *http.Request) {
	newProduceRequest := &[]models.Produce{}
	err := json.NewDecoder(r.Body).Decode(newProduceRequest) //decode the request body into struct and failed if any error occur
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validProduce, invalidProduce := s.filterNewProduceRequest(*newProduceRequest)
	createdProduce, failedProduce := s.createAllProduce(validProduce)

	failedProduce = append(failedProduce, invalidProduce...)

	js, err := json.Marshal(&models.CreateProduceResponse{
		Created: createdProduce,
		Invalid: failedProduce,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status := http.StatusCreated
	switch {
	case len(failedProduce) > 0 && len(createdProduce) > 0:
		status = http.StatusMultiStatus
	case len(failedProduce) > 0 && len(createdProduce) == 0:
		status = http.StatusBadRequest
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return
}

func (s *Server) filterNewProduceRequest(newProduce []models.Produce) ([]models.Produce, []models.Produce) {
	validProduce := []models.Produce{}
	failedProduce := []models.Produce{}
	for _, val := range newProduce {
		if _, exists := s.data[strings.ToUpper(val.ProduceCode)]; exists {
			failedProduce = append(failedProduce, val)
			continue
		}

		if !isValidProduceCode(val.ProduceCode) {
			failedProduce = append(failedProduce, val)
			continue
		}

		validProduce = append(validProduce, val)
	}
	return validProduce, failedProduce
}

func isValidProduceCode(produceCode string) bool {
	validProduceCode := regexp.MustCompile(`[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{4}`)
	return validProduceCode.MatchString(produceCode)
}

// createAllProduce will attempt to add every produce passed in by newProduce
// It will return a slice of all created produce and failed produce for error reporting.
func (s *Server) createAllProduce(newProduce []models.Produce) ([]models.Produce, []models.Produce) {
	type createResponse struct {
		successful *models.Produce
		failed     *models.Produce
	}

	createdProduce := []models.Produce{}
	failedProduce := []models.Produce{}
	ch := make(chan createResponse)
	wgSelect := sync.WaitGroup{}
	wgSelect.Add(1)

	go func() {
		for ch != nil {
			select {
			case x, ok := <-ch:
				if !ok {
					ch = nil
					wgSelect.Done()
					break
				}
				if x.successful != nil {
					s.mutex.Lock()
					createdProduce = append(createdProduce, *x.successful)
					s.mutex.Unlock()
				}
				if x.failed != nil {
					s.mutex.Lock()
					failedProduce = append(failedProduce, *x.successful)
					s.mutex.Unlock()
				}
			}
		}
	}()

	wg := sync.WaitGroup{}
	for _, val := range newProduce {
		wg.Add(1)
		go func(val models.Produce, ch chan createResponse) {
			defer wg.Done()

			p, err := s.createProduce(val)
			if err != nil {
				ch <- createResponse{
					failed: &val,
				}
				return
			}

			ch <- createResponse{
				successful: &p,
			}
		}(val, ch)
	}

	wg.Wait()
	close(ch)
	wgSelect.Wait()

	return createdProduce, failedProduce
}

// createProduce cannot return an error in this implementation
// however this method is part of an interface so a future implementation that uses
// a real database would need the option to return an error
func (s *Server) createProduce(produce models.Produce) (models.Produce, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	newProduce := models.Produce{
		Name:        produce.Name,
		ProduceCode: strings.ToUpper(produce.ProduceCode),
		UnitPrice:   math.Round(produce.UnitPrice*100) / 100,
	}

	s.data[newProduce.ProduceCode] = newProduce

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
