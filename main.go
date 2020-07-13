package main

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"guru/handlers"
	"guru/repositories"
	"guru/services"
	"log"
	"net/http"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	logger, _ := zap.NewProduction()
	mode, exists := os.LookupEnv("MODE")
	if exists && mode == "development" {
		logger, _ = zap.NewDevelopment()
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error(err.Error())
		}
	}()

	mongoUser, exist := os.LookupEnv("MONGO_INITDB_ROOT_USERNAME")
	if !exist {
		mongoUser = "mongo"
	}

	mongoPassword, exist := os.LookupEnv("MONGO_INITDB_ROOT_PASSWORD")
	if !exist {
		mongoPassword = "mongo"
	}

	mongoHost, exist := os.LookupEnv("MONGO_HOST")
	if !exist {
		mongoHost = "localhost"
	}

	mongoPort, exist := os.LookupEnv("MONGO_PORT")
	if !exist {
		mongoPort = "27017"
	}

	mongoDb, exist := os.LookupEnv("MONGO_INITDB_DATABASE")
	if !exist {
		mongoDb = "guru"
	}

	credential := options.Credential{
		Username: mongoUser,
		Password: mongoPassword,
	}
	clientOpts := options.Client().ApplyURI("mongodb://" + mongoHost + ":" + mongoPort).SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		logger.Fatal(err.Error())
	}

	db := client.Database(mongoDb)

	service := &services.UserService{
		UserRepository:        repositories.UserRepository{DB: db},
		DepositRepository:     repositories.DepositRepository{DB: db},
		TransactionRepository: repositories.TransactionRepository{DB: db},
	}
	service.Run()

	r := router{
		userHandler:        handlers.NewUserHandler(service),
		transactionHandler: handlers.NewTransactionHandler(service),
	}

	log.Fatal(http.ListenAndServe(":8080", r.InitRouter()))
}
