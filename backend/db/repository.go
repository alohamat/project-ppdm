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

	// construtor
func NewUserRepository(col *mongo.Collection) *UserRepository {
	return &UserRepository{
		collection: col,
	}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) (bson.ObjectID, error) {
	res, err := r.collection.InsertOne(ctx, u)
	if err != nil {
		return bson.NilObjectID, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	id, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		return bson.NilObjectID, errors.New("não foi possível obter o objectid inserido")
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
		return nil, fmt.Errorf("erro ao buscar usuário por id: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar usuário por email: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar usuários: %w", err)
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("erro ao decodificar usuários: %w", err)
	}
	if users == nil {
		users = []models.User{}
	}
	return users, nil
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("erro ao contar usuários: %w", err)
	}
	return count, nil
}

func (r *UserRepository) UpdateRole(ctx context.Context, id bson.ObjectID, newRole string) error {
	res, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"role": newRole}},
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar cargo: %w", err)
	}
	if res.MatchedCount == 0 {
		return errors.New("usuário não encontrado")
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}
	if res.DeletedCount == 0 {
		return errors.New("usuário não encontrado")
	}
	return nil
}
