package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/xmattstrongx/supermarket/controllers"

	"github.com/gorilla/mux"
)

func Serve() {
	router := mux.NewRouter()

	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8080" //localhost
	}

	router.HandleFunc("/api/user/new", controllers.CreateProduce).Methods("POST")

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
