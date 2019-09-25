package helpers

import (
	"recibe_me/configs"

	"gopkg.in/mgo.v2"
)

// IssuesCollection is a IssuesCollection
var IssuesCollection = getSession().DB(configs.DefaultServerConfig.MongoDbName).C(configs.DefaultServerConfig.MongoCollectionIssues)

// DeliveriesCollection is a DeliveriesCollection
var DeliveriesCollection = getSession().DB(configs.DefaultServerConfig.MongoDbName).C(configs.DefaultServerConfig.MongoCollectionDeliveries)

// UsersCollection is a UsersCollection
var UsersCollection = getSession().DB(configs.DefaultServerConfig.MongoDbName).C(configs.DefaultServerConfig.MongoCollectionUsers)

func getSession() *mgo.Session {
	session, err := mgo.Dial(configs.DefaultServerConfig.MongoURL)

	if err != nil {
		panic(err)
	}

	return session
}
