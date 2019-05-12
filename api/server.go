package api

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/xmattstrongx/supermarket/models"
)

// ProduceManager is the interface between the API and the backend storage
type ProduceManager interface {
	ListProduce(queryParameters) ([]models.Produce, error)
	CreateProduce(models.Produce) (models.Produce, error)
	DeleteProduce(string) error
}

// Server is the data structure for holding all types needed by the server to run and serve requests
type Server struct {
	produceManager ProduceManager
	backend
	// data    map[string]models.Produce
}

type backend struct {
	data  map[string]models.Produce
	once  sync.Once
	mutex sync.RWMutex
}

// NewServer instantiates a new Server
func NewServer() *Server {
	server := &Server{}
	server.backend.once.Do(func() {
		server.backend.data = initializeData()
	})
	return server
}

// Serve starts a Server for handling API requests
func (s *Server) Serve() {
	router := mux.NewRouter()

	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8000" //localhost
	}

	router.HandleFunc("/api/v1/produce", s.ListProduce).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/produce", s.CreateProduce).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/produce/{productCode}", s.DeleteProduce).Methods(http.MethodDelete)

	log.Infof("Listening and serving on :%s", port)
	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}

func initializeData() map[string]models.Produce {
	data := map[string]models.Produce{
		"A12T-4GH7-QPL9-3N4M": models.Produce{
			Name:        "Lettuce",
			ProduceCode: "A12T-4GH7-QPL9-3N4M",
			UnitPrice:   3.46,
		},
		"E5T6-9UI3-TH15-QR88": models.Produce{
			Name:        "Peach",
			ProduceCode: "E5T6-9UI3-TH15-QR88",
			UnitPrice:   2.99,
		},
		"YRT6-72AS-K736-L4AR": models.Produce{
			Name:        "Green Pepper",
			ProduceCode: "YRT6-72AS-K736-L4AR",
			UnitPrice:   0.79,
		},
		"TQ4C-VV6T-75ZX-1RMR": models.Produce{
			Name:        "Gala Apple",
			ProduceCode: "TQ4C-VV6T-75ZX-1RMR",
			UnitPrice:   3.59,
		},
	}
	return data
}
