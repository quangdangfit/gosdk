package dbs

import (
	"gopkg.in/mgo.v2"
	"transport/lib/utils/logger"
)

type MongoDB struct {
	*mgo.Database
}

func (db *MongoDB) EnsureIndex(collectionName string, index mgo.Index) bool {
	sessionClone := db.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.Name).C(collectionName)

	err := collection.EnsureIndex(index)
	if err != nil {
		logger.Error("[EnsureIndex] Ensure index name fail: ", index.Name, " - ERROR: ", err)
		return false
	}
	logger.Info("[EnsureIndex] Successful to ensure index: ", index.Name)
	return true
}

func (db *MongoDB) DropIndex(collectionName string, name string) bool {
	sessionClone := db.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.Name).C(collectionName)

	err := collection.DropIndexName(name)
	if err != nil {
		logger.Error("[DropIndex] Drop index fail: ", name, " - Error: ", err)
		return false
	}
	logger.Info("[DropIndex] Drop index success: ", name)
	return true
}

// -- Func find one record
func (db *MongoDB) FindOne(collectionName string, query interface{}, sort string, TResult interface{}) (err error) {
	sessionClone := db.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.Name).C(collectionName)

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

// -- Func find many record
func (db *MongoDB) FindMany(collectionName string, query map[string]interface{}, sort string, TResult interface{}) (err error) {
	sessionClone := db.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.Name).C(collectionName)

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

// -- Func find many record
func (db *MongoDB) FindManyPaging(collectionName string, query map[string]interface{}, sort string, page int, limit int, TResult interface{}) (*Paging, error) {
	sessionClone := db.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.Name).C(collectionName)

	cursor := collection.Find(query).Sort(sort)
	total, _ := cursor.Count()
	if total == 0 {
		return nil, mgo.ErrNotFound
	}
	paging := GetPaging(page, limit, total)
	err := cursor.Skip(paging.Skip).Limit(paging.Limit).All(TResult)
	if err == mgo.ErrNotFound {
		logger.Error("[FindManyPaging] ERROR: ", err)
		return nil, err
	}

	return &paging, nil
}

// -- Func pipe all
func (db *MongoDB) PipeAll(collectionName string, pipeline interface{}, TResult interface{}) (err error) {
	sessionClone := db.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.Name).C(collectionName)

	err = collection.Pipe(pipeline).All(TResult)
	if err != nil {
		logger.Error("[PipeAll] ERROR: ", err)
		return err
	}

	logger.Info("[InsertOne] Successful to pipe")
	return nil
}

// -- Func Insert one record
func (db *MongoDB) InsertOne(collectionName string, payload interface{}) (err error) {
	sessionClone := db.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.Name).C(collectionName)

	err = collection.Insert(payload)
	if err != nil {
		logger.Error("[InsertOne] ERROR: ", err)
		return err
	}

	logger.Info("[InsertOne] Successful to insert record")
	return nil
}

// -- Func Insert many records
func (db *MongoDB) InsertMany(collectionName string, payload []interface{}) (err error) {
	sessionClone := db.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.Name).C(collectionName)

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

// -- Func Upsert record
func (db *MongoDB) Upsert(collectionName string, selector interface{}, payload interface{}) (err error) {
	sessionClone := db.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.Name).C(collectionName)

	_, err = collection.Upsert(selector, payload)
	if err != nil {
		logger.Error("[Upsert] ERROR: ", err)
		return err
	}

	logger.Info("[Upsert] Successful to upsert record")
	return nil
}

// -- Func Update one record
func (db *MongoDB) UpdateOne(collectionName string, selector interface{}, payload interface{}) (err error) {
	sessionClone := db.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.Name).C(collectionName)

	err = collection.Update(selector, payload)
	if err != nil {
		logger.Error("[UpdateOne] ERROR: ", err)
		return err
	}

	logger.Info("[UpdateOne] Successful to update the record")
	return nil
}

// -- Func Update many records
func (db *MongoDB) UpdateMany(collectionName string, selector interface{}, payload interface{}) (err error) {
	sessionClone := db.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.Name).C(collectionName)

	changed, err := collection.UpdateAll(selector, payload)
	if err != nil {
		logger.Error("[UpdateMany] ERROR: ", err)
		return err
	}

	logger.Infof("[UpdateMany] Matched %d and update %d records", changed.Matched, changed.Updated)
	return nil
}

// -- Func Delete one record
func (db *MongoDB) DeleteOne(collectionName string, selector interface{}) (err error) {
	sessionClone := db.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.Name).C(collectionName)

	err = collection.Remove(selector)
	if err != nil {
		logger.Error("[DeleteOne] ERROR: ", err)
		return err
	}

	logger.Info("[DeleteOne] Successful to delete the record")
	return nil
}

// -- Func Delete many record
func (db *MongoDB) DeleteMany(collectionName string, selector interface{}) (err error) {
	sessionClone := db.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.Name).C(collectionName)

	changed, err := collection.RemoveAll(selector)
	if err != nil {
		logger.Error("[DeleteMany] ERROR: ", err)
		return err
	}

	logger.Infof("[DeleteMany] Matched %d and deleted %d records", changed.Matched, changed.Removed)
	return nil
}

// -- Func Apply record
func (db *MongoDB) ApplyDB(collectionName string, selector interface{}, payload interface{}, TResult interface{}) (err error) {
	sessionClone := db.Session.Copy()
	defer sessionClone.Close()
	collection := sessionClone.DB(db.Name).C(collectionName)

	change := mgo.Change{Update: payload, ReturnNew: true}
	_, err = collection.Find(selector).Apply(change, TResult)
	if err != nil {
		logger.Error("[ApplyDB] ERROR: ", err)
		return err
	}

	logger.Info("[ApplyDB] Successful to apply")
	return nil
}
