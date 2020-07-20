package main

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	db "github.com/quangdangfit/gosdk/database"
	"github.com/quangdangfit/gosdk/database/mongo"
)

var Database mongo.Database

type Brand struct {
	Code string `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}

func index() {
	index := mgo.Index{
		Key:        []string{"code"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     false,
		Name:       "index_brand_code",
	}
	Database.EnsureIndex("brand", index)
	Database.DropIndex("brand", "index_brand_code")
}

func AddBrand() {
	payload := map[string]string{
		"code": "code",
		"name": "name",
	}
	Database.InsertOne("brand", payload)
}

func AddManyBrand() {
	payload := map[string]string{
		"code": "code1",
		"name": "name1",
	}

	payload1 := map[string]string{
		"code": "code2",
		"name": "brand2",
	}
	listPayload := []interface{}{payload, payload1}
	Database.InsertMany("brand", listPayload)
}

func UpdateBrand() {
	filter := bson.M{"code": "code"}
	update := bson.M{"$set": map[string]string{"code": "quang1"}}
	Database.UpdateOne("brand", filter, update)
}

func UpdateManyBrand() {
	filter := bson.M{"code": "code"}
	update := bson.M{"$set": map[string]string{"code": "quang1"}}
	Database.UpdateMany("brand", filter, update)
}

func DeleteBrand() {
	filter := bson.M{"code": "quang"}
	Database.DeleteOne("brand", filter)
	Database.DeleteMany("brand", filter)
}

func init() {
	dbConfig := db.Config{
		Hosts:        "localhost:27017",
		AuthDatabase: "admin",
		AuthUserName: "",
		AuthPassword: "",
		Database:     "testdb",
	}

	Database = mongo.New(dbConfig)
}

func main() {
	dbConfig := db.Config{
		Hosts:        "localhost:27017",
		AuthDatabase: "admin",
		AuthUserName: "",
		AuthPassword: "",
		Database:     "testdb",
	}

	Database = mongo.New(dbConfig)

	var results = []Brand{}

	filter := bson.M{"code": "code"}
	Database.FindMany("brand", filter, "-_id", &results)

	for _, e := range results {
		log.Println(e.Name, e.Code)
	}
}
