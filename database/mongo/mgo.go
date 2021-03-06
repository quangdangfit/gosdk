package mongo

import (
	"log"
	"strings"
	"time"

	"gopkg.in/mgo.v2"

	db "github.com/quangdangfit/gosdk/database"
	"github.com/quangdangfit/gosdk/utils/logger"
	"github.com/quangdangfit/gosdk/utils/paging"
)

type gomgo struct {
	conn *mgo.Database
}

func NewMongo(config db.Config) db.Mongo {
	var s *mgo.Session

	if config.Env == "replica" {
		s = replicaSetConnection(config)
	} else {
		s = nativeConnection(config)
	}
	db := mgo.Database{Session: s, Name: config.Database}
	return &gomgo{conn: &db}
}

func NewWithConnString(uri string) db.Mongo {
	dialInfo, err := mgo.ParseURL(uri)
	if err != nil {
		logger.Fatal("Connection string is invalid")
	}

	dialInfo.Timeout = time.Duration(10)
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		logger.Fatal(err)
	}

	db := mgo.Database{Session: session, Name: dialInfo.Database}
	return &gomgo{conn: &db}
}

func nativeConnection(config db.Config) *mgo.Session {
	logger.Info("[nativeConnection] Connecting mongodb")

	var timeout = 10
	if config.ConnectionTimeout > 0 {
		timeout = config.ConnectionTimeout
	}
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{config.Hosts},
		Timeout:  time.Duration(timeout) * time.Second,
		Database: config.AuthDatabase,
		Username: config.AuthUserName,
		Password: config.AuthPassword,
	}
	// Create a session which maintains a pool of socket connections
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		logger.Fatal("[nativeConnection] Failed to connect gomgo: ", err)
	}

	session.SetSafe(&mgo.Safe{})
	logger.Info("[nativeConnection] Mongodb connected")
	return session
}

func replicaSetConnection(config db.Config) *mgo.Session {
	logger.Info("[replicaSetConnection] Connecting mongodb")

	var timeout = 10
	if config.ConnectionTimeout > 0 {
		timeout = config.ConnectionTimeout
	}
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:          strings.Split(config.Hosts, ","),
		Username:       config.AuthUserName,
		Password:       config.AuthPassword,
		Database:       config.AuthDatabase,
		ReplicaSetName: config.Replica,
		Timeout:        time.Duration(timeout) * time.Second,
	})
	if err != nil {
		log.Println("[replicaSetConnection] ERROR: ", err)
	}

	session.SetSafe(&mgo.Safe{})
	logger.Info("[replicaSetConnection] Connected to replica set success!")

	return session
}

// GetSession will return a copy session from session of database, remember close session at the end.
func (db *gomgo) GetSession() *mgo.Session {
	return db.conn.Session.Copy()
}

// EnsureIndex ensures an index with the given collection name and key exists, creating it with
// the provided parameters if necessary. EnsureIndex does not modify a previously
// existent index with a matching key. The old index must be dropped first instead.
func (db *gomgo) EnsureIndex(collectionName string, index mgo.Index) bool {
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
func (db *gomgo) DropIndex(collectionName string, name string) bool {
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
func (db *gomgo) FindOne(collectionName string, query map[string]interface{}, sort string, TResult interface{}) (err error) {
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
func (db *gomgo) FindMany(collectionName string, query map[string]interface{}, sort string, TResult interface{}) (err error) {
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
func (db *gomgo) FindManyPaging(collectionName string, query map[string]interface{}, sort string, page int, limit int, TResult interface{}) (*paging.Paging, error) {
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
func (db *gomgo) PipeAll(collectionName string, pipeline interface{}, TResult interface{}) (err error) {
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
func (db *gomgo) InsertOne(collectionName string, payload interface{}) (err error) {
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
func (db *gomgo) InsertMany(collectionName string, payload []interface{}) (err error) {
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
func (db *gomgo) Upsert(collectionName string, selector interface{}, payload interface{}) (err error) {
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
func (db *gomgo) UpdateOne(collectionName string, selector interface{}, payload interface{}) (err error) {
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
func (db *gomgo) UpdateMany(collectionName string, selector interface{}, payload interface{}) (err error) {
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
func (db *gomgo) DeleteOne(collectionName string, selector interface{}) (err error) {
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
func (db *gomgo) DeleteMany(collectionName string, selector interface{}) (err error) {
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
func (db *gomgo) ApplyDB(collectionName string, selector interface{}, payload interface{}, TResult interface{}) (err error) {
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
