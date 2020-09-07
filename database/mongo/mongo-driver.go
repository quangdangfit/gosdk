package mongo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	db "github.com/quangdangfit/gosdk/database"
	"github.com/quangdangfit/gosdk/utils/logger"
	"github.com/quangdangfit/gosdk/utils/paging"
)

type mongodb struct {
	conn *mongo.Database
}

func NewWithConfig(config db.Config) *mongodb {
	logger.Info("Connecting mongodb, database: ", config.Database)

	connectionURI := "mongodb://"
	if config.AuthUserName != "" && config.AuthPassword != "" {
		connectionURI += fmt.Sprintf("%s:%s@", config.AuthUserName, config.AuthPassword)
	}
	connectionURI += config.Hosts
	if config.AuthDatabase != "" {
		connectionURI += fmt.Sprintf("/?authSource=%s", config.AuthDatabase)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))

	if err != nil {
		logger.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Mongodb connected")
	return &mongodb{conn: client.Database(config.Database)}
}

func NewConnection(uri string) *mongodb {
	dbname := ""
	temp := strings.Split(uri, "/")
	if len(temp) == 4 {
		dbname = strings.Split(temp[3], "?")[0]
	}

	logger.Info("Connecting mongodb, database: ", dbname)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		logger.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Mongodb connected")
	return &mongodb{conn: client.Database(dbname)}
}

// FindOne executes the query and unmarshals the first obtained document into the
// result argument, using the provided collection name, query and sort.
// The query may be a map or a struct value capable of being marshalled with bson.
// The sort is field name need to sort, a field name may be prefixed by - (minus) for
// it to be sorted in reverse order.
func (db *mongodb) FindOne(collectionName string, filter map[string]interface{}, sort interface{}, result interface{}) (err error) {
	collection := db.conn.Collection(collectionName)
	ctx := context.TODO()

	opts := options.FindOne()
	if sort != nil {
		opts.SetSort(sort)
	}

	err = collection.FindOne(ctx, filter).Decode(result)
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
func (db *mongodb) FindMany(collectionName string, filter map[string]interface{}, sort interface{}, results interface{}) (err error) {
	collection := db.conn.Collection(collectionName)
	ctx := context.TODO()

	opts := options.Find()
	if sort != nil {
		opts.SetSort(sort)
	}

	cur, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return err
	}
	err = cur.All(ctx, results)
	if err != nil {
		return err
	}

	if err := cur.Err(); err != nil {
		return err
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	return nil
}

// FindMany executes FindMany function but skip by page parameter and limit by limit parameter
func (db *mongodb) FindManyPaging(collectionName string, filter map[string]interface{}, sort interface{}, page int, limit int, results interface{}) (*paging.Paging, error) {
	collection := db.conn.Collection(collectionName)

	opts := options.Find()
	if sort != nil {
		opts.SetSort(sort)
	}
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64((page - 1) * limit))

	ctx := context.TODO()
	total, err := collection.CountDocuments(ctx, filter)
	pagingObj := paging.New(page, limit, int(total))

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, results)
	if err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	return pagingObj, nil
}

// Insert inserts one document in the respective collection, the returned error will
// be an error.
func (db *mongodb) InsertOne(collectionName string, payload interface{}) (err error) {
	collection := db.conn.Collection(collectionName)
	ctx := context.TODO()

	_, err = collection.InsertOne(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

// InsertMany queues up the provided documents for insertion and run insert.
func (db *mongodb) InsertMany(collectionName string, payload []interface{}) (err error) {
	collection := db.conn.Collection(collectionName)
	ctx := context.TODO()

	_, err = collection.InsertMany(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

// UpdateOne finds a single document matching the provided selector document
// and modifies it according to the update document.
func (db *mongodb) UpdateOne(collectionName string, filter interface{}, payload interface{}) (err error) {
	collection := db.conn.Collection(collectionName)
	ctx := context.TODO()

	_, err = collection.UpdateOne(ctx, filter, payload)
	if err != nil {
		return err
	}

	return nil
}

// UpdateMany finds all documents matching the provided selector document
// and modifies them according to the update document.
func (db *mongodb) UpdateMany(collectionName string, selector interface{}, payload interface{}) (err error) {
	collection := db.conn.Collection(collectionName)
	ctx := context.TODO()

	_, err = collection.UpdateMany(ctx, selector, payload)
	if err != nil {
		return err
	}

	return nil
}

// DeleteOne finds a single document matching the provided selector document
// and removes it from the database.
func (db *mongodb) DeleteOne(collectionName string, filter interface{}) (err error) {
	collection := db.conn.Collection(collectionName)
	ctx := context.TODO()

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMany finds all documents matching the provided selector document
// and removes them from the database.
func (db *mongodb) DeleteMany(collectionName string, selector interface{}) (err error) {
	collection := db.conn.Collection(collectionName)
	ctx := context.TODO()

	_, err = collection.DeleteMany(ctx, selector)
	if err != nil {
		return err
	}

	return nil
}
