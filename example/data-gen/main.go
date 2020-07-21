package main

import (
	"context"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/xo/dburl"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// This script is to add random data to the db
// to demonstrate the output/

func main() {
	time.Sleep(time.Duration(10) * time.Second)
	ctx := context.Background()
	// postgres connection
	postgresconnectionEnv := "POSTGRES_CONN"
	postgresconnectionStr := os.Getenv(postgresconnectionEnv)
	db, err := dburl.Open(postgresconnectionStr)
	if err != nil {
		log.Fatal(err)
	}

	// mongo connection
	mongoconnectionEnv := "MONGO_CONN"
	mongoconnectionString := os.Getenv(mongoconnectionEnv)
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoconnectionString))
	if err != nil {
		log.Fatal(err)
	}
	err = mongoClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	dropTableStmt := `DROP TABLE IF EXISTS templates`
	_, err = db.Exec(dropTableStmt)
	if err != nil {
		log.Fatal(err)
	}
	// create table in postgres
	createTableStmt := `CREATE TABLE templates(name TEXT, is_active BOOLEAN)`
	_, err = db.Exec(createTableStmt)
	if err != nil {
		log.Fatal(err)
	}

	// Insert rows
	insertStmt := `INSERT INTO templates (name, is_active) VALUES ($1, $2)`
	_, err = db.Exec(insertStmt, "test1", true)
	if err != nil {
		log.Fatal(err)
	}
	// Insert documents in mongo
	coll := mongoClient.Database("test").Collection("test")
	_, err = coll.InsertOne(ctx, map[string]interface{}{"name": "test", "is_active": true})
	if err != nil {
		log.Fatal(err)
	}
	_, err = coll.InsertOne(ctx, map[string]interface{}{"name": "test1", "is_active": false})
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Duration(3) * time.Second)
	_, err = coll.InsertOne(ctx, map[string]interface{}{"name": "test2", "is_active": false})
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(insertStmt, "test2", true)
	if err != nil {
		log.Fatal(err)
	}
}
