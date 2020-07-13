package main

import (
	"github.com/gorilla/mux"
	"guru/handlers"
	"net/http"
)

type router struct {
	userHandler *handlers.UserHandler
	transactionHandler *handlers.TransactionHandler
}

func (router router) InitRouter() *mux.Router {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./swaggerui/"))
	r.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))

	s := r.PathPrefix("/user").Subrouter()
	s.HandleFunc("/create", router.userHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/get", router.userHandler.Get).Methods(http.MethodPost)
	s.HandleFunc("/deposit", router.userHandler.AddDeposit).Methods(http.MethodPost)

	r.HandleFunc("/transaction", router.transactionHandler.Transaction).Methods(http.MethodPost)

	return r
}
