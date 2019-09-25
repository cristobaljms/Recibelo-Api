package configs

// ServerConfig server cofiguration
type ServerConfig struct {
	Port                      string
	MongoURL                  string
	MongoDbName               string
	MongoCollectionUsers      string
	MongoCollectionDeliveries string
	MongoCollectionIssues     string
}

// DefaultServerConfig server cofiguration by default
var DefaultServerConfig = ServerConfig{
	Port:                      "0.0.0.0:8070",
	MongoURL:                  "mongodb://192.168.1.192:27017",
	MongoDbName:               "recibeme_db",
	MongoCollectionUsers:      "users",
	MongoCollectionDeliveries: "deliveries",
	MongoCollectionIssues:     "issues",
}
