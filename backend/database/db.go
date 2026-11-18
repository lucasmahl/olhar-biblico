package database

import (
	"context"
	"fmt"
	"log"
	"olhar-biblico/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Conn() (*mongo.Client, error) {
	ctx := context.Background()

	str := fmt.Sprintf("mongodb://%s:%s@localhost:%s", config.DbUser, config.DbPassword, config.DbPort)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(str))
	if err != nil {
		return nil, err
	}

	// Testa a conexão.
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	return client, nil
}
