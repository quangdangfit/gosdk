package mongo

import (
	"log"
	"strings"
	"time"

	"gopkg.in/mgo.v2"

	db "github.com/quangdangfit/gosdk/database"
)

func nativeConnection(config db.Config) *mgo.Session {
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
		log.Println("[nativeConnection] ERROR: ", err)
	}
	log.Println("[nativeConnection] Mongodb connected")
	return mgoSession
}

func replicaSetConnection(config db.Config) *mgo.Session {
	log.Println("[replicaSetConnection] Connecting mongodb")

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:          strings.Split(config.Hosts, ","),
		Username:       config.AuthUserName,
		Password:       config.AuthPassword,
		Database:       config.AuthDatabase,
		ReplicaSetName: config.Replica,
	})
	if err != nil {
		log.Println("[replicaSetConnection] ERROR: ", err)
	}

	session.SetSafe(&mgo.Safe{})
	log.Println("[replicaSetConnection] Connected to replica set success!")

	return session
}
