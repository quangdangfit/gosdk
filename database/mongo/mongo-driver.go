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

	opts := options.FindOne()
	if sort != nil {
		opts.SetSort(sort)
	}

	err = collection.FindOne(context.TODO(), query).Decode(result)
	if err != nil {
		return err
	}

	return nil
}
