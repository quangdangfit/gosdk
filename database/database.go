package database

import (
	"strings"
	"time"
	"transport/lib/utils/logger"

	"gopkg.in/mgo.v2"
)

type DBConfig struct {
	MongoDBHosts      string
	AuthDatabase      string
	Database          string
	AuthUserName      string
	AuthPassword      string
	ConnectionTimeout int64
	Env               string
	Replica           string
}

func NewConnection(config DBConfig) MongoDB {
	var s *mgo.Session

	if config.Env == "replica" {
		s = replicaSetConnection(config)
	} else {
		s = nativeConnection(config)
	}
	db := mgo.Database{Session: s, Name: config.Database}
	return &mongoDB{conn: &db}
}

func nativeConnection(config DBConfig) *mgo.Session {
	logger.Info("[nativeConnection] Connecting mongodb")

	var timeout time.Duration = 60
	if config.ConnectionTimeout > 0 {
		timeout = time.Duration(config.ConnectionTimeout)
	}
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{config.MongoDBHosts},
		Timeout:  timeout * time.Second,
		Database: config.AuthDatabase,
		Username: config.AuthUserName,
		Password: config.AuthPassword,
	}
	// Create a session which maintains a pool of socket connections
	mgoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		logger.Fatal("[nativeConnection] ERROR: ", err)
	}
	logger.Info("[nativeConnection] Mongodb connected")
	return mgoSession
}

func replicaSetConnection(config DBConfig) *mgo.Session {
	logger.Info("[replicaSetConnection] Connecting mongodb")

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:          strings.Split(config.MongoDBHosts, ","),
		Username:       config.AuthUserName,
		Password:       config.AuthPassword,
		Database:       config.AuthDatabase,
		ReplicaSetName: config.Replica,
	})
	if err != nil {
		logger.Error("[replicaSetConnection] ERROR: ", err)
		panic(err)
	}

	session.SetSafe(&mgo.Safe{})
	logger.Info("[replicaSetConnection] Connected to replica set success!")

	return session
}