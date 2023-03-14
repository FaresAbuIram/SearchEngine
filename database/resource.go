package database

import (
	"context"
	"mime/multipart"
	"searchEngine/models"

	"gopkg.in/mgo.v2/bson"
)

func InsertNewResource(ctx context.Context, reqResource models.ReqResource, file *multipart.FileHeader, fileName string) error {

	data := bson.M{
		"title": reqResource.Title,
		"type": reqResource.Type,
		"tags": reqResource.Tags,
		"path": fileName + file.Filename,
	}
	_, err := SearchCollection.InsertOne(Ctx, data)

	return err
}
