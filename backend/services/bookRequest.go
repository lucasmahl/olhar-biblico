package services

import (
	"olhar-biblico/config"
	"olhar-biblico/database"
	"olhar-biblico/models"
	"olhar-biblico/models/search_filter"
	"olhar-biblico/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookRequestService interface {
	FindByEmailOrPhone(email, phone string) (models.BookRequest, error)
	Create(bookRequest *models.BookRequest) (interface{}, error)
	Find(filter search_filter.SearchFilter, bookRequestFilter models.BookRequest) (utils.Result, error)
}

type bookRequestService struct{}

func NewBookRequestService() BookRequestService {
	return &bookRequestService{}
}

func (s *bookRequestService) FindByEmailOrPhone(email, phone string) (models.BookRequest, error) {
	ctx, cancel := config.GetContext()
	defer cancel()

	client, err := database.Conn()
	if err != nil {
		return models.BookRequest{}, err
	}
	defer client.Disconnect(ctx)

	filter := bson.D{{Key: "$or",
		Value: []bson.D{
			{{Key: "email", Value: email}},
			{{Key: "phone", Value: phone}},
		},
	}}

	bookRequest := models.BookRequest{}
	collection := client.Database(config.DbName).Collection(bookRequest.TableName())

	collection.FindOne(ctx, filter).Decode(&bookRequest)

	return bookRequest, nil
}

func (s *bookRequestService) Create(bookRequest *models.BookRequest) (interface{}, error) {
	ctx, cancel := config.GetContext()
	defer cancel()

	client, err := database.Conn()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	collection := client.Database(config.DbName).Collection(bookRequest.TableName())

	now := time.Now()
	bookRequest.CreatedAt = &now
	res, err := collection.InsertOne(ctx, &bookRequest)
	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

func (s *bookRequestService) Find(filter search_filter.SearchFilter, bookRequestFilter models.BookRequest) (utils.Result, error) {
	result := utils.Result{}
	ctx, cancel := config.GetContext()
	defer cancel()

	client, err := database.Conn()
	if err != nil {
		return result, err
	}
	defer client.Disconnect(ctx)

	pipeline := bson.D{}
	if bookRequestFilter.Name != "" {
		pipeline = append(
			pipeline,
			bson.E{Key: "name", Value: bson.D{{Key: "$regex", Value: bookRequestFilter.Name}, {Key: "$options", Value: "i"}}},
		)
	}
	if bookRequestFilter.Email != "" {
		pipeline = append(
			pipeline,
			bson.E{Key: "email", Value: bson.D{{Key: "$regex", Value: bookRequestFilter.Email}, {Key: "$options", Value: "i"}}},
		)
	}
	if bookRequestFilter.Address != "" {
		pipeline = append(
			pipeline,
			bson.E{Key: "address", Value: bson.D{{Key: "$regex", Value: bookRequestFilter.Address}, {Key: "$options", Value: "i"}}},
		)
	}
	if bookRequestFilter.Phone != "" {
		pipeline = append(
			pipeline,
			bson.E{Key: "phone", Value: bson.D{{Key: "$regex", Value: bookRequestFilter.Phone}, {Key: "$options", Value: "i"}}},
		)
	}

	collection := client.Database(config.DbName).Collection(bookRequestFilter.TableName())

	// Busca para contar a quantidade de objetos.
	if result.TotalRecord, err = collection.CountDocuments(ctx, pipeline); err != nil {
		return result, err
	}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: filter.Order}})
	opts.SetSkip(int64(filter.Skip))
	opts.SetLimit(int64(filter.Limit))

	cur, err := collection.Find(ctx, pipeline, opts)
	if err != nil {
		return result, err
	}
	defer cur.Close(ctx)

	bookRequest := []models.BookRequest{}
	if err := cur.All(ctx, &bookRequest); err != nil {
		return result, err
	}

	result.Results = bookRequest
	result.Limit = uint64(filter.Limit)
	result.Skip = uint64(filter.Skip)
	result.Paginable = true
	if filter.Skip == 0 && int64(filter.Limit) >= result.TotalRecord {
		result.Paginable = false
	}

	return result, nil
}
