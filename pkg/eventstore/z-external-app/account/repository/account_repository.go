package repository

import (
	"context"
	"eventstore-intro/pkg/eventstore/z-external-app/account/models"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AccountMongoRepository interface {
	//Insert(ctx context.Context, order *models.AccountProjection) (string, error)
	GetByID(ctx context.Context, orderID string) (*models.AccountProjection, error)
	UpdateAccount(ctx context.Context, order *models.AccountProjection) error
	DeactivateAccount(ctx context.Context, order *models.AccountProjection) error
	ActivateAccount(ctx context.Context, order *models.AccountProjection) (string, error)
}

func (m *MongoRepository) ActivateAccount(ctx context.Context, account *models.AccountProjection) (string, error) {
	_, err := m.db.Database("emoney").Collection("accounts").InsertOne(ctx, account, &options.InsertOneOptions{})
	if err != nil {
		return "", err
	}
	return account.AccountNumber, nil
}
