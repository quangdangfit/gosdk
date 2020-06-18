package database

import (
	"gitlab.com/quangdangfit/gocommon/utils/paging"
	"gopkg.in/mgo.v2"
)

//TODO remove log

type Database interface {
	GetSession() *mgo.Session
	EnsureIndex(collectionName string, index mgo.Index) bool
	DropIndex(collectionName string, name string) bool
	FindOne(collectionName string, query interface{}, sort string, TResult interface{}) (err error)
	FindMany(collectionName string, query map[string]interface{}, sort string, TResult interface{}) (err error)
	FindManyPaging(collectionName string, query map[string]interface{}, sort string, page int, limit int, TResult interface{}) (*paging.Paging, error)
	PipeAll(collectionName string, pipeline interface{}, TResult interface{}) (err error)
	InsertOne(collectionName string, payload interface{}) (err error)
	InsertMany(collectionName string, payload []interface{}) (err error)
	Upsert(collectionName string, selector interface{}, payload interface{}) (err error)
	UpdateOne(collectionName string, selector interface{}, payload interface{}) (err error)
	UpdateMany(collectionName string, selector interface{}, payload interface{}) (err error)
	DeleteOne(collectionName string, selector interface{}) (err error)
	DeleteMany(collectionName string, selector interface{}) (err error)
	ApplyDB(collectionName string, selector interface{}, payload interface{}, TResult interface{}) (err error)
}

type database struct{}
