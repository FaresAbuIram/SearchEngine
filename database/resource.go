package database

import (
	"context"
	"log"
	"mime/multipart"
	"searchEngine/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func InsertNewResource(ctx context.Context, reqResource models.ReqResource, file *multipart.FileHeader, fileName string) error {

	data := bson.M{
		"title": reqResource.Title,
		"type":  reqResource.Type,
		"tags":  reqResource.Tags,
		"path":  fileName + file.Filename,
	}
	_, err := SearchCollection.InsertOne(Ctx, data)

	return err
}

func GetResources(ctx context.Context, keyword string) ([]models.SearchEngineResult, error) {
	filter := bson.M{
		"tags": bson.M{"$regex": keyword},
	}

	var result []models.SearchEngineResult
	resources, err := SearchCollection.Find(ctx, filter)
	if err != nil {
		log.Println("Error get data from mongo: ", err)
		return nil, err
	}
	resources.All(ctx, &result)
	return result, nil
}

func GetResource(ctx context.Context, id string) (models.Resource, error) {

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
		log.Println("Error get data from mongo: ", err)
		return models.Resource{}, err
	}
	
	return result, nil
}
