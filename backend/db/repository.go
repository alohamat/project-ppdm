package db

import (
	"context"
	"errors"
	"fmt"
	"ppdm/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

// construtor da struct
func NewUserRepository(col *mongo.Collection) *UserRepository {
	return &UserRepository{
		collection: col,
	}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) (bson.ObjectID, error) {
	res, err := r.collection.InsertOne(ctx, u)
	if err != nil {
		return bson.NilObjectID, fmt.Errorf("erro ao criar user: %w", err)
	}

	id, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		return bson.NilObjectID, errors.New("não foi possível obter objectid inserido")
	}

	u.ID = id
	return id, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id bson.ObjectID) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar usuario por id: %w", err)
	}
	return &user, nil
}