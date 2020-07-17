package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"guru/models"
	"guru/repositories"
	"guru/services"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

var srv *httptest.Server

func TestMain(m *testing.M) {
	credential := options.Credential{
		Username: "mongo",
		Password: "mongo",
	}
	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err.Error())
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-stop
		log.Printf("[WARN] interrupt signal")
		cancel()
	}()

	db := client.Database("test_guru")
	service := &services.UserService{
		UserRepository:        repositories.UserRepository{DB: db},
		DepositRepository:     repositories.DepositRepository{DB: db},
		TransactionRepository: repositories.TransactionRepository{DB: db},
		Ticker:                time.NewTicker(10 * time.Second),
	}


	userHandler := NewUserHandler(service)
	transactionHandler := NewTransactionHandler(service)

	r := mux.NewRouter()

	s := r.PathPrefix("/user").Subrouter()
	s.HandleFunc("/create", userHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/get", userHandler.Get).Methods(http.MethodPost)
	s.HandleFunc("/deposit", userHandler.AddDeposit).Methods(http.MethodPost)

	r.HandleFunc("/transaction", transactionHandler.Transaction).Methods(http.MethodPost)

	service.Run(ctx, r)

	srv = httptest.NewServer(r)
	defer srv.Close()

	os.Exit(m.Run())
}

func TestTransactionHandler_Transaction(t *testing.T) {
	jsonStr := []byte(`{
 		"user_id": 1,
		"transaction_id": 1,
		"type": "Win",
		"amount": 25,
		"token": "sssss"
	}`)

	res, err := http.Post(
		fmt.Sprintf("%s/transaction", srv.URL),
		"application/json",
		bytes.NewBuffer(jsonStr),
	)

	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var transactionResponse models.TransactionResponseModel
	if err = json.Unmarshal(resBytes, &transactionResponse); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, models.TransactionResponseModel{Error: "", Balance: 75}, transactionResponse)
}

func TestTransactionHandler_TransactionWrongToken(t *testing.T) {
	jsonStr := []byte(`{
 		"user_id": 1,
		"transaction_id": 1,
		"type": "Win",
		"amount": 25,
		"token": "ttttt"
	}`)

	res, err := http.Post(
		fmt.Sprintf("%s/transaction", srv.URL),
		"application/json",
		bytes.NewBuffer(jsonStr),
	)

	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var errorResponse models.ErrorResponseModel
	if err = json.Unmarshal(resBytes, &errorResponse); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, models.ErrorResponseModel{Error: "wrong token"}, errorResponse)
}

func TestTransactionHandler_TransactionNotFound(t *testing.T) {
	jsonStr := []byte(`{
 		"user_id": 5,
		"transaction_id": 1,
		"type": "Win",
		"amount": 25,
		"token": "ttttt"
	}`)

	res, err := http.Post(
		fmt.Sprintf("%s/transaction", srv.URL),
		"application/json",
		bytes.NewBuffer(jsonStr),
	)

	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var errorResponse models.ErrorResponseModel
	if err = json.Unmarshal(resBytes, &errorResponse); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, models.ErrorResponseModel{Error: "not found"}, errorResponse)
}

func TestTransactionHandler_TransactionNotEnoughBalance(t *testing.T) {
	jsonStr := []byte(`{
 		"user_id": 1,
		"transaction_id": 1,
		"type": "Bet",
		"amount": 300,
		"token": "sssss"
	}`)

	res, err := http.Post(
		fmt.Sprintf("%s/transaction", srv.URL),
		"application/json",
		bytes.NewBuffer(jsonStr),
	)

	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var errorResponse models.ErrorResponseModel
	if err = json.Unmarshal(resBytes, &errorResponse); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, models.ErrorResponseModel{Error: "not enough balance"}, errorResponse)
}
