package database

import (
	"gitlab.com/quangdangfit/gocommon/utils/logger"
	"gitlab.com/quangdangfit/gocommon/utils/paging"

	"gopkg.in/mgo.v2"
)

type mongoDB struct {
	conn *mgo.Database
}

type MongoDB interface {
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

// GetSession will return a copy session from session of database, remember close session at the end.
func (db *mongoDB) GetSession() *mgo.Session {
	return db.conn.Session.Copy()
}

// EnsureIndex ensures an index with the given collection name and key exists, creating it with
// the provided parameters if necessary. EnsureIndex does not modify a previously
// existent index with a matching key. The old index must be dropped first instead.
func (db *mongoDB) EnsureIndex(collectionName string, index mgo.Index) bool {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	err := collection.EnsureIndex(index)
	if err != nil {
		logger.Error("[EnsureIndex] Ensure index name fail: ", index.Name, " - ERROR: ", err)
		return false
	}
	logger.Info("[EnsureIndex] Successful to ensure index: ", index.Name)
	return true
}

// DropIndex removes the index with the provided collection name and index name.
func (db *mongoDB) DropIndex(collectionName string, name string) bool {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	err := collection.DropIndexName(name)
	if err != nil {
		logger.Error("[DropIndex] Drop index fail: ", name, " - Error: ", err)
		return false
	}
	logger.Info("[DropIndex] Drop index success: ", name)
	return true
}

// FindOne executes the query and unmarshals the first obtained document into the
// result argument, using the provided collection name, query and sort.
// The query may be a map or a struct value capable of being marshalled with bson.
// The sort is field name need to sort, a field name may be prefixed by - (minus) for
// it to be sorted in reverse order.
func (db *mongoDB) FindOne(collectionName string, query interface{}, sort string, TResult interface{}) (err error) {
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
		logger.Error("[FindOne] ERROR: ", err)
		return err
	}

	return nil
}

// FindMany executes the query and unmarshals the all obtained document into the
// result argument, using the provided collection name, query, sort, and result interface.
// The query may be a map or a struct value capable of being marshalled with bson.
// The sort is field name need to sort, a field name may be prefixed by - (minus) for
// it to be sorted in reverse order.
func (db *mongoDB) FindMany(collectionName string, query map[string]interface{}, sort string, TResult interface{}) (err error) {
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
		logger.Error("[FindMany] ERROR: ", err)
		return err
	}

	return nil
}

// FindMany executes FindMany function but skip by page parameter and limit by limit parameter
func (db *mongoDB) FindManyPaging(collectionName string, query map[string]interface{}, sort string, page int, limit int, TResult interface{}) (*paging.Paging, error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	cursor := collection.Find(query).Sort(sort)
	total, _ := cursor.Count()
	if total == 0 {
		return nil, mgo.ErrNotFound
	}
	pagingObj := paging.GetPaging(page, limit, total)
	err := cursor.Skip(pagingObj.Skip).Limit(pagingObj.Limit).All(TResult)
	if err == mgo.ErrNotFound {
		logger.Error("[FindManyPaging] ERROR: ", err)
		return nil, err
	}

	return &pagingObj, nil
}

// PipeAll prepares a pipeline to aggregate. The pipeline document
// must be a slice built in terms of the aggregation framework language.
func (db *mongoDB) PipeAll(collectionName string, pipeline interface{}, TResult interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	err = collection.Pipe(pipeline).All(TResult)
	if err != nil {
		logger.Error("[PipeAll] ERROR: ", err)
		return err
	}

	logger.Info("[InsertOne] Successful to pipe")
	return nil
}

// Insert inserts one document in the respective collection, the returned error will
// be an error.
func (db *mongoDB) InsertOne(collectionName string, payload interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	err = collection.Insert(payload)
	if err != nil {
		logger.Error("[InsertOne] ERROR: ", err)
		return err
	}

	logger.Info("[InsertOne] Successful to insert record")
	return nil
}

// InsertMany queues up the provided documents for insertion and run insert.
func (db *mongoDB) InsertMany(collectionName string, payload []interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	bulk := collection.Bulk()
	bulk.Insert(payload...)
	_, err = bulk.Run()

	if err != nil {
		logger.Error("[InsertMany] ERROR: ", err)
		return err
	}
	logger.Info("[InsertMany] Successful to inserted records")
	return nil
}

// Upsert finds a single document matching the provided selector document
// and modifies it according to the update document.
func (db *mongoDB) Upsert(collectionName string, selector interface{}, payload interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	_, err = collection.Upsert(selector, payload)
	if err != nil {
		logger.Error("[Upsert] ERROR: ", err)
		return err
	}

	logger.Info("[Upsert] Successful to upsert record")
	return nil
}

// UpdateOne finds a single document matching the provided selector document
// and modifies it according to the update document.
func (db *mongoDB) UpdateOne(collectionName string, selector interface{}, payload interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	err = collection.Update(selector, payload)
	if err != nil {
		logger.Error("[UpdateOne] ERROR: ", err)
		return err
	}

	logger.Info("[UpdateOne] Successful to update the record")
	return nil
}

// UpdateMany finds all documents matching the provided selector document
// and modifies them according to the update document.
func (db *mongoDB) UpdateMany(collectionName string, selector interface{}, payload interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	changed, err := collection.UpdateAll(selector, payload)
	if err != nil {
		logger.Error("[UpdateMany] ERROR: ", err)
		return err
	}

	logger.Infof("[UpdateMany] Matched %d and update %d records", changed.Matched, changed.Updated)
	return nil
}

// DeleteOne finds a single document matching the provided selector document
// and removes it from the database.
func (db *mongoDB) DeleteOne(collectionName string, selector interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	err = collection.Remove(selector)
	if err != nil {
		logger.Error("[DeleteOne] ERROR: ", err)
		return err
	}

	logger.Info("[DeleteOne] Successful to delete the record")
	return nil
}

// DeleteMany finds all documents matching the provided selector document
// and removes them from the database.
func (db *mongoDB) DeleteMany(collectionName string, selector interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	changed, err := collection.RemoveAll(selector)
	if err != nil {
		logger.Error("[DeleteMany] ERROR: ", err)
		return err
	}

	logger.Infof("[DeleteMany] Matched %d and deleted %d records", changed.Matched, changed.Removed)
	return nil
}

// Apply runs the findAndModify MongoDB command, which allows updating, upserting
// or removing a document matching a query and atomically returning either the old
// version (the default) or the new version of the document.
func (db *mongoDB) ApplyDB(collectionName string, selector interface{}, payload interface{}, TResult interface{}) (err error) {
	sessionClone := db.conn.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.conn.Name).C(collectionName)

	change := mgo.Change{Update: payload, ReturnNew: true}
	_, err = collection.Find(selector).Apply(change, TResult)
	if err != nil {
		logger.Error("[ApplyDB] ERROR: ", err)
		return err
	}

	logger.Info("[ApplyDB] Successful to apply")
	return nil
}
