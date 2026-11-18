package services

import (
	"olhar-biblico/config"
	"olhar-biblico/database"
	"olhar-biblico/models"
	"olhar-biblico/models/search_filter"
	"olhar-biblico/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type VisitRequestService interface {
	FindByEmailOrPhone(email, phone string) (models.VisitRequest, error)
	Create(visitRequest *models.VisitRequest) (interface{}, error)
	Find(filter search_filter.SearchFilter, visitRequestFilter models.VisitRequest) (utils.Result, error)
}

type visitRequestService struct{}

func NewVisitRequestService() VisitRequestService {
	return &visitRequestService{}
}

func (s *visitRequestService) FindByEmailOrPhone(email, phone string) (models.VisitRequest, error) {
	ctx, cancel := config.GetContext()
	defer cancel()

	client, err := database.Conn()
	if err != nil {
		return models.VisitRequest{}, err
	}
	defer client.Disconnect(ctx)

	filter := bson.D{{Key: "$or",
		Value: []bson.D{
			{{Key: "email", Value: email}},
			{{Key: "phone", Value: phone}},
		},
	}}

	visitRequest := models.VisitRequest{}
	collection := client.Database(config.DbName).Collection(visitRequest.TableName())

	collection.FindOne(ctx, filter).Decode(&visitRequest)

	return visitRequest, nil
}

func (s *visitRequestService) Create(visitRequest *models.VisitRequest) (interface{}, error) {
	ctx, cancel := config.GetContext()
	defer cancel()

	client, err := database.Conn()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	collection := client.Database(config.DbName).Collection(visitRequest.TableName())

	now := time.Now()
	visitRequest.CreatedAt = &now
	res, err := collection.InsertOne(ctx, &visitRequest)
	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

func (s *visitRequestService) Find(filter search_filter.SearchFilter, visitRequestFilter models.VisitRequest) (utils.Result, error) {
	result := utils.Result{}
	ctx, cancel := config.GetContext()
	defer cancel()

	client, err := database.Conn()
	if err != nil {
		return result, err
	}
	defer client.Disconnect(ctx)

	var pipeline mongo.Pipeline = nil
	if visitRequestFilter.Name != "" {
		pipeline = append(
			pipeline,
			bson.D{{
				Key:   "$match",
				Value: bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: visitRequestFilter.Name}, {Key: "$options", Value: "i"}}}},
			}},
		)
	}
	if visitRequestFilter.Email != "" {
		pipeline = append(
			pipeline,
			bson.D{{
				Key:   "$match",
				Value: bson.D{{Key: "email", Value: bson.D{{Key: "$regex", Value: visitRequestFilter.Email}, {Key: "$options", Value: "i"}}}},
			}},
		)
	}
	if visitRequestFilter.Address != "" {
		pipeline = append(
			pipeline,
			bson.D{{
				Key:   "$match",
				Value: bson.D{{Key: "address", Value: bson.D{{Key: "$regex", Value: visitRequestFilter.Address}, {Key: "$options", Value: "i"}}}},
			}},
		)
	}
	if visitRequestFilter.Phone != "" {
		pipeline = append(
			pipeline,
			bson.D{{
				Key:   "$match",
				Value: bson.D{{Key: "phone", Value: bson.D{{Key: "$regex", Value: visitRequestFilter.Phone}, {Key: "$options", Value: "i"}}}},
			}},
		)
	}

	collection := client.Database(config.DbName).Collection(visitRequestFilter.TableName())

	// Busca para contar a quantidade de objetos.
	if cur, err := collection.Aggregate(ctx, pipeline); err != nil {
		return result, err
	} else {
		defer cur.Close(ctx)

		// Para apagar os valores das chaves, que são sobrescritos, na consulta acima é recriada a variavel na linha abaixo
		// que é repetida nas proximas linhas
		visitRequest := []models.VisitRequest{}
		if err := cur.All(ctx, &visitRequest); err != nil {
			return result, err
		}
		result.TotalRecord = int64(len(visitRequest))
	}

	pipeline = append(
		pipeline,
		bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: filter.Order}}}},
		bson.D{{Key: "$skip", Value: filter.Skip}},
		bson.D{{Key: "$limit", Value: filter.Limit}},
	)

	cur, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return result, err
	}
	defer cur.Close(ctx)

	visitRequest := []models.VisitRequest{}
	if err := cur.All(ctx, &visitRequest); err != nil {
		return result, err
	}

	result.Results = visitRequest
	result.Limit = uint64(filter.Limit)
	result.Skip = uint64(filter.Skip)
	result.Paginable = true
	if filter.Skip == 0 && int64(filter.Limit) >= result.TotalRecord {
		result.Paginable = false
	}

	return result, nil
}
