package main

import (
	"crud/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/user", service.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/user", service.FindAllUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}", service.FindOneUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}", service.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/user/{id}", service.DeleteUser).Methods(http.MethodDelete)

	fmt.Println("Escutando na porta 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
