package main

import (
	"gopkg.in/mgo.v2/bson"
	"log"
)

func main() {
	var results = []Brand{}

	filter := bson.M{"code": "code"}
	Database.FindMany("brand", filter, "-_id", &results)

	for _, e := range results {
		log.Println(e.Name, e.Code)
	}
}
