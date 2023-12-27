package repository

import "go.mongodb.org/mongo-driver/mongo"

type MongoRepository struct {
	db *mongo.Client
}

func NewMongoRepository(db *mongo.Client) *MongoRepository {
	return &MongoRepository{db: db}
}
