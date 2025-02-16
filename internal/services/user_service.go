package services

import (
	"context"
	"errors"

	"github.com/joaogustavosp/loyalty-api/internal/models"
	"github.com/joaogustavosp/loyalty-api/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	collection *mongo.Collection
}

func NewUserService(collection *mongo.Collection) *UserService {
	return &UserService{collection}
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) (*mongo.InsertOneResult, error) {
	if !utils.IsValidCPF(user.CPF) {
		return nil, errors.New("invalid CPF")
	}
	exists, err := s.UserExistsByCPF(ctx, user.CPF)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("CPF already exists")
	}
	return s.collection.InsertOne(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, id string) (models.User, error) {
	var user models.User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}
	filter := bson.M{"_id": objID}
	err = s.collection.FindOne(ctx, filter).Decode(&user)
	return user, err
}

func (s *UserService) UpdateUser(ctx context.Context, id string, user models.User) (*mongo.UpdateResult, error) {
	if !utils.IsValidCPF(user.CPF) {
		return nil, errors.New("CPF inválido")
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Encontrar o usuário atual
	var existingUser models.User
	filter := bson.M{"_id": objID}
	err = s.collection.FindOne(ctx, filter).Decode(&existingUser)
	if err != nil {
		return nil, err
	}

	// Manter a senha existente
	if user.Password == "" {
		user.Password = existingUser.Password
	}

	// Manter o status existente se não for fornecido
	if user.Status == "" {
		user.Status = existingUser.Status
	}
	update := bson.M{"$set": user}
	return s.collection.UpdateOne(ctx, filter, update)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	return s.collection.DeleteOne(ctx, filter)
}

func (s *UserService) UserExistsByCPF(ctx context.Context, cpf string) (bool, error) {
	filter := bson.M{"cpf": cpf}
	count, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, userID, hashedPassword string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"password": hashedPassword, "status": models.Active}}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	return err
}

func (s *UserService) GetUserByCPF(ctx context.Context, cpf string) (models.User, error) {
	var user models.User
	filter := bson.M{"cpf": cpf}
	err := s.collection.FindOne(ctx, filter).Decode(&user)
	return user, err
}
