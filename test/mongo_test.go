package test

import (
	"GoIm/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestFindOne(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		Username: "admin",
		Password: "admin",
	}).ApplyURI("mongodb://150.158.144.5:27017"))
	if err != nil {
		t.Fatal(err)
	}
	db := client.Database("im")
	ub := new(models.UserBasic)
	err = db.Collection("user_basic").FindOne(context.Background(), bson.D{{"account", "moyuyiren"}}).Decode(&ub)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("ub====>", ub)
}

func TestFind(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		Username: "admin",
		Password: "admin",
	}).ApplyURI("mongodb://150.158.144.5:27017"))
	if err != nil {
		t.Fatal(err)
	}
	db := client.Database("im")
	cursor, err := db.Collection("user_room").
		Find(context.Background(), bson.D{})
	if err != nil {
		t.Fatal(err)
	}
	urs := make([]*models.UserRoom, 0)
	for cursor.Next(context.Background()) {
		ub := new(models.UserRoom)
		err := cursor.Decode(&ub)
		if err != nil {
			t.Fatal(err)
		}
		urs = append(urs, ub)
	}

	for _, ur := range urs {
		fmt.Println("UserRoom====>", ur)
	}
}
