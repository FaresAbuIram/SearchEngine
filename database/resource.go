package database

import (
	"context"
	"log"
	"mime/multipart"
	"searchEngine/models"

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
	filter := bson.M{
		"_id": bson.ObjectIdHex(id),
	}

	var result models.Resource
	resource, err := SearchCollection.Find(ctx, filter)
	if err != nil {
		log.Println("Error get data from mongo: ", err)
		return result, err
	}

	resource.All(ctx, &result)
	return result, nil
}
