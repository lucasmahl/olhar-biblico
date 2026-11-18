package services

import (
	"olhar-biblico/config"
	"olhar-biblico/database"
	"olhar-biblico/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService interface {
	Create(user *models.User) (interface{}, error)
	FindByEmail(email string) (models.User, error)
	UpdatePassword(user *models.User) (interface{}, error)
	FindUserWithoutPasswordByID(id string) (models.User, error)
}

type userService struct{}

func NewUserService() UserService {
	return &userService{}
}

func (s *userService) Create(user *models.User) (interface{}, error) {
	ctx, cancel := config.GetContext()
	defer cancel()

	client, err := database.Conn()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	collection := client.Database(config.DbName).Collection(user.TableName())

	now := time.Now()
	user.CreatedAt = &now
	res, err := collection.InsertOne(ctx, &user)
	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

func (s *userService) FindByEmail(email string) (models.User, error) {
	user := models.User{}

	ctx, cancel := config.GetContext()
	defer cancel()

	client, err := database.Conn()
	if err != nil {
		return user, err
	}
	defer client.Disconnect(ctx)

	collection := client.Database(config.DbName).Collection(user.TableName())
	collection.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&user)

	return user, nil
}

func (s *userService) FindUserWithoutPasswordByID(id string) (models.User, error) {
	userWithoutPassword := models.User{}

	ctx, cancel := config.GetContext()
	defer cancel()

	client, err := database.Conn()
	if err != nil {
		return userWithoutPassword, err
	}
	defer client.Disconnect(ctx)

	collection := client.Database(config.DbName).Collection(userWithoutPassword.TableName())

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return userWithoutPassword, err
	}

	projection := options.FindOne().SetProjection(bson.D{{Key: "password", Value: 0}})

	collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}, projection).Decode(&userWithoutPassword)

	return userWithoutPassword, nil
}

func (s *userService) UpdatePassword(user *models.User) (interface{}, error) {
	ctx, cancel := config.GetContext()
	defer cancel()

	client, err := database.Conn()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	collection := client.Database(config.DbName).Collection(user.TableName())

	now := time.Now()
	user.UpdatedAt = &now
	update := bson.M{
		"$set": bson.M{
			"password": user.Password,
		},
	}

	res, err := collection.UpdateByID(ctx, user.ID, update)
	if err != nil {
		return nil, err
	}

	return res.UpsertedID, nil
}
