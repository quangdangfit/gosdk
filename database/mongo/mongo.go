package mongo

import (
	"log"

	"gopkg.in/mgo.v2"

	db "github.com/quangdangfit/gocommon/database"
	"github.com/quangdangfit/gocommon/utils/paging"
)

type Database interface {
	GetSession() *mgo.Session
	EnsureIndex(collectionName string, index mgo.Index) bool
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

type mongodb struct {
	conn *mgo.Database
}

func New(config db.Config) Database {
	var s *mgo.Session

	if config.Env == "replica" {
		s = replicaSetConnection(config)
	} else {
		s = nativeConnection(config)
	}
	db := mgo.Database{Session: s, Name: config.Database}
	return &mongodb{conn: &db}
}

// GetSession will return a copy session from session of database, remember close session at the end.
func (db *mongodb) GetSession() *mgo.Session {
	return db.conn.Session.Copy()
}

// EnsureIndex ensures an index with the given collection name and key exists, creating it with
// the provided parameters if necessary. EnsureIndex does not modify a previously
// existent index with a matching key. The old index must be dropped first instead.
func (db *mongodb) EnsureIndex(collectionName string, index mgo.Index) bool {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	err := collection.EnsureIndex(index)
	if err != nil {
		log.Println("[EnsureIndex] Ensure index name fail: ", index.Name, " - ERROR: ", err)
		return false
	}
	log.Println("[EnsureIndex] Successful to ensure index: ", index.Name)
	return true
}

// DropIndex removes the index with the provided collection name and index name.
func (db *mongodb) DropIndex(collectionName string, name string) bool {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	err := collection.DropIndexName(name)
	if err != nil {
		log.Println("[DropIndex] Drop index fail: ", name, " - Error: ", err)
		return false
	}
	log.Println("[DropIndex] Drop index success: ", name)
	return true
}

// FindOne executes the query and unmarshals the first obtained document into the
// result argument, using the provided collection name, query and sort.
// The query may be a map or a struct value capable of being marshalled with bson.
// The sort is field name need to sort, a field name may be prefixed by - (minus) for
// it to be sorted in reverse order.
func (db *mongodb) FindOne(collectionName string, query map[string]interface{}, sort string, TResult interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	if sort != "" {
		cursor := collection.Find(query).Sort(sort)
		err = cursor.One(TResult)
	} else {
		err = collection.Find(query).One(TResult)
	}

	if err != nil {
		return err
	}

	return nil
}

// FindMany executes the query and unmarshals the all obtained document into the
// result argument, using the provided collection name, query, sort, and result interface.
// The query may be a map or a struct value capable of being marshalled with bson.
// The sort is field name need to sort, a field name may be prefixed by - (minus) for
// it to be sorted in reverse order.
func (db *mongodb) FindMany(collectionName string, query map[string]interface{}, sort string, TResult interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	if sort != "" {
		cursor := collection.Find(query).Sort(sort)
		err = cursor.All(TResult)
	} else {
		err = collection.Find(query).All(TResult)
	}

	if err != nil {
		return err
	}

	return nil
}

// FindMany executes FindMany function but skip by page parameter and limit by limit parameter
func (db *mongodb) FindManyPaging(collectionName string, query map[string]interface{}, sort string, page int, limit int, TResult interface{}) (*paging.Paging, error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	cursor := collection.Find(query).Sort(sort)
	total, _ := cursor.Count()
	if total == 0 {
		return nil, mgo.ErrNotFound
	}
	pagingObj := paging.New(page, limit, total)
	err := cursor.Skip(pagingObj.Skip).Limit(pagingObj.Limit).All(TResult)
	if err == mgo.ErrNotFound {
		return nil, err
	}

	return pagingObj, nil
}

// PipeAll prepares a pipeline to aggregate. The pipeline document
// must be a slice built in terms of the aggregation framework language.
func (db *mongodb) PipeAll(collectionName string, pipeline interface{}, TResult interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	err = collection.Pipe(pipeline).All(TResult)
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts one document in the respective collection, the returned error will
// be an error.
func (db *mongodb) InsertOne(collectionName string, payload interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	err = collection.Insert(payload)
	if err != nil {
		return err
	}

	return nil
}

// InsertMany queues up the provided documents for insertion and run insert.
func (db *mongodb) InsertMany(collectionName string, payload []interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	bulk := collection.Bulk()
	bulk.Insert(payload...)
	_, err = bulk.Run()

	if err != nil {
		return err
	}
	return nil
}

// Upsert finds a single document matching the provided selector document
// and modifies it according to the update document.
func (db *mongodb) Upsert(collectionName string, selector interface{}, payload interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	_, err = collection.Upsert(selector, payload)
	if err != nil {
		return err
	}

	return nil
}

// UpdateOne finds a single document matching the provided selector document
// and modifies it according to the update document.
func (db *mongodb) UpdateOne(collectionName string, selector interface{}, payload interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	err = collection.Update(selector, payload)
	if err != nil {
		return err
	}

	return nil
}

// UpdateMany finds all documents matching the provided selector document
// and modifies them according to the update document.
func (db *mongodb) UpdateMany(collectionName string, selector interface{}, payload interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	_, err = collection.UpdateAll(selector, payload)
	if err != nil {
		return err
	}

	return nil
}

// DeleteOne finds a single document matching the provided selector document
// and removes it from the database.
func (db *mongodb) DeleteOne(collectionName string, selector interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	err = collection.Remove(selector)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMany finds all documents matching the provided selector document
// and removes them from the database.
func (db *mongodb) DeleteMany(collectionName string, selector interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	_, err = collection.RemoveAll(selector)
	if err != nil {
		return err
	}

	return nil
}

// Apply runs the findAndModify Database command, which allows updating, upserting
// or removing a document matching a query and atomically returning either the old
// version (the default) or the new version of the document.
func (db *mongodb) ApplyDB(collectionName string, selector interface{}, payload interface{}, TResult interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	change := mgo.Change{Update: payload, ReturnNew: true}
	_, err = collection.Find(selector).Apply(change, TResult)
	if err != nil {
		return err
	}

	return nil
}
