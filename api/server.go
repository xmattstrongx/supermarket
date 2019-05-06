package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/xmattstrongx/supermarket/models"
)

type Server struct {
	data map[string]models.Produce
}

func NewServer() *Server {
	server := &Server{
		data: initializeData(),
	}
	return server
}

func (s *Server) Serve() {
	router := mux.NewRouter()

	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8080" //localhost
	}

	router.HandleFunc("/api/v1/produce", s.listProduce).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/produce", s.createProduce).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/produce/{productCode}", s.deleteProduce).Methods(http.MethodDelete)

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
