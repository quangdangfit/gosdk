package mongo

import (
	db "gitlab.com/quangdangfit/gocommon/database"
	"gopkg.in/mgo.v2"
	"log"
	"strings"
	"time"
)

type MongoDB interface {
	db.Database
}

func New(config db.DBConfig) MongoDB {
	var s *mgo.Session

	if config.Env == "replica" {
		s = replicaSetConnection(config)
	} else {
		s = nativeConnection(config)
	}
	db := mgo.Database{Session: s, Name: config.Database}
	return &mongoDB{conn: &db}
}

func nativeConnection(config db.DBConfig) *mgo.Session {
	log.Println("[nativeConnection] Connecting mongodb")

	var timeout time.Duration = 60
	if config.ConnectionTimeout > 0 {
		timeout = time.Duration(config.ConnectionTimeout)
	}
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{config.Hosts},
		Timeout:  timeout * time.Second,
		Database: config.AuthDatabase,
		Username: config.AuthUserName,
		Password: config.AuthPassword,
	}
	// Create a session which maintains a pool of socket connections
	mgoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalln("[nativeConnection] ERROR: ", err)
	}
	log.Println("[nativeConnection] Mongodb connected")
	return mgoSession
}

func replicaSetConnection(config db.DBConfig) *mgo.Session {
	log.Println("[replicaSetConnection] Connecting mongodb")

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:          strings.Split(config.Hosts, ","),
		Username:       config.AuthUserName,
		Password:       config.AuthPassword,
		Database:       config.AuthDatabase,
		ReplicaSetName: config.Replica,
	})
	if err != nil {
		log.Fatalln("[replicaSetConnection] ERROR: ", err)
	}

	session.SetSafe(&mgo.Safe{})
	log.Println("[replicaSetConnection] Connected to replica set success!")

	return session
}
