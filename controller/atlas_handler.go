package controller

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connect() (*mongo.Client, context.Context) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://admin:urUf3IZyhtTqFEuE@tutuplapak-cluster.xkljd.mongodb.net/tutuplapak"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client.Connect(ctx)

	return client, ctx
}
