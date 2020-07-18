package database

import (
	"gitlab.com/quangdangfit/gocommon/utils/paging"
)

type Database interface {
	EnsureIndex(collectionName string, index IndexConfig) bool
	DropIndex(collectionName string, name string) bool
	FindOne(table string, query map[string]interface{}, sort string, result interface{}) (err error)
	FindMany(table string, query map[string]interface{}, sort string, result interface{}) (err error)
	FindManyPaging(table string, query map[string]interface{}, sort string, offset int, limit int, result interface{}) (*paging.Paging, error)
	PipeAll(table string, pipeline interface{}, result interface{}) (err error)
	InsertOne(table string, payload interface{}) (err error)
	InsertMany(table string, payload []interface{}) (err error)
	Upsert(table string, selector interface{}, payload interface{}) (err error)
	UpdateOne(table string, selector interface{}, payload interface{}) (err error)
	UpdateMany(table string, selector interface{}, payload interface{}) (err error)
	DeleteOne(table string, selector interface{}) (err error)
	DeleteMany(table string, selector interface{}) (err error)
	ApplyDB(table string, selector interface{}, payload interface{}, result interface{}) (err error)
}
