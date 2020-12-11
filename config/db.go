package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"Prodkt/controllers"
)

func Connect() {
	// Database Config
	//mongodb+srv://<username>:<password>@cluster0-zzart.mongodb.net/test?retryWrites=true&w=majority
	//mongodb://Imran:12345@localhost:27017/?authSource=admin
	//clientOptions := options.Client().ApplyURI("mongodb://imran:12345@localhost:27017/test?authSource=admin&replicaSet=Cluster0-shard-0&readPreference=primary&appname=MongoDB%20Compass&ssl=true")
	clientOptions := options.Client().ApplyURI("mongodb://Imran:12345@localhost:27017/?authSource=admin")
	client, err := mongo.NewClient(clientOptions)

	//Set up a context required by mongo.Connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	//To close the connection at the end
	defer cancel()

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
	}
	db := client.Database("Imran")
	controllers.DeviceCollection(db)
	controllers.ServiceCollection(db)
	controllers.UserCollection(db)
	return
}