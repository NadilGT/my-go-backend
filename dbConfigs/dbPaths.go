package dbConfigs

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	DATABASE *mongo.Database
	CLIENT   *mongo.Client
)

const DATABASE_URL = "mongodb+srv://admin:W6ptbj7HPS3RJ4cU@cluster0.tgypip5.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
const DATABASE_NAME = "POS"
