package database

import (
	"context"
	"searchEngine/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type Resource interface {
	InsertNewResource(ctx context.Context, reqResource models.ReqResource, fileName, filePath string) error
	GetResources(ctx context.Context, keyword string) ([]models.SearchEngineResult, error)
	GetResource(ctx context.Context, id string) (models.Resource, error)
}

type ResourceDatabase struct{}

func NewResourceService() *ResourceDatabase {
	return &ResourceDatabase{}
}

// InsertNewResource function to insert new resource in the database.
func (r *ResourceDatabase) InsertNewResource(ctx context.Context, reqResource models.ReqResource, fileName, filePath string) error {
	data := bson.M{
		"title": reqResource.Title,
		"type":  reqResource.Type,
		"tags":  reqResource.Tags,
		"path":  filePath + fileName,
	}
	_, err := SearchCollection.InsertOne(Ctx, data)
	if err != nil {
		return err
	}
	return nil
}

// GetResources function to get all the resources
// that partially or fully match the keyword string.
func (r *ResourceDatabase) GetResources(ctx context.Context, keyword string) ([]models.SearchEngineResult, error) {
	filter := bson.M{
		"tags": bson.M{"$regex": keyword},
	}

	var result []models.SearchEngineResult
	resources, err := SearchCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	resources.All(ctx, &result)
	return result, nil
}

// GetResource function to get resource's information by id.
func (r *ResourceDatabase) GetResource(ctx context.Context, id string) (models.Resource, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Resource{}, err
	}
	filter := bson.M{
		"_id": objectId,
	}

	var result models.Resource
	err = SearchCollection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		return models.Resource{}, err
	}

	return result, nil
}
