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

func (db *mongodb) FindOne(collectionName string, query map[string]interface{}, sort interface{}, result interface{}) (err error) {
	collection := db.conn.Collection(collectionName)
	ctx := context.TODO()

	opts := options.FindOne()
	if sort != nil {
		opts.SetSort(sort)
	}

	err = collection.FindOne(ctx, query).Decode(result)
	if err != nil {
		return err
	}

	return nil
}

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

func (db *mongodb) InsertOne(collectionName string, payload interface{}) (err error) {
	collection := db.conn.Collection(collectionName)
	ctx := context.TODO()

	_, err = collection.InsertOne(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}
