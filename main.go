package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburi  = "mongodb://localhost:27017"
	dbname = "xmlToJsonApi"
)

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dburi))
	if err != nil {
		panic(err)
	}
}
