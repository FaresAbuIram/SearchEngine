package controllers

import (
	"log"
	"net/http"

	"searchEngine/database"
	"searchEngine/models"

	"github.com/gin-gonic/gin"
)

func ErrorResponseStatus(context *gin.Context, status int, message string) {
	context.JSON(status, gin.H{"message": message})
}

func SuccessResponseStatus(context *gin.Context, status int, data string) {
	context.JSON(status, gin.H{"message": data})

}

func SuccessResponseStatusForToken(context *gin.Context, status int, data interface{}) {
	cacceptedHeader := context.Request.Header["Accept"][0]
	switch cacceptedHeader {
	case "text/xml":
		context.XML(status, gin.H{"data": data, "message": "success"})
	default:
		context.JSON(status, gin.H{"data": data, "message": "success"})
	}
}

// Create new resource
// @Summary      Create a new resource
// @Description  This route uses to create a new resource
// @Accept       mpfd
// @Produce      json
// @Param 		 file  formData file true "File"
// @Param        data  formData models.ReqResource true "Create new resource"
// @Success      200  {object}  models.CreateResourceResponse
// @Failure      400  {object}	models.CreateResourceResponse
// @Failure      500  {object}	models.CreateResourceResponse
// @Router       /createNewResource [post]
func Createresource(context *gin.Context) {
	var reqResource models.ReqResource

	if context.PostForm("title") != "" {
		reqResource.Title = context.PostForm("title")
	} else {
		log.Printf("title is required")
		ErrorResponseStatus(context, http.StatusBadRequest, "title is required")
		return
	}

	if context.PostForm("type") != "" {
		reqResource.Type = context.PostForm("type")
	} else {
		log.Printf("type is required")
		ErrorResponseStatus(context, http.StatusBadRequest, "type is required")
		return
	}

	if len(context.PostFormArray("tags")) != 0 {
		reqResource.Tags = context.PostFormArray("tags")
	} else {
		log.Printf("tags is required")
		ErrorResponseStatus(context, http.StatusBadRequest, "tags is required")
		return
	}

	file, err := context.FormFile("file")
	if err != nil {
		log.Printf("File Error: ", err)
		ErrorResponseStatus(context, http.StatusBadRequest, "Can't upload the file")
		return
	}

	var fileName string

	if reqResource.Type == "HTML Pages" {
		fileName = "files/HTML/"

	} else if reqResource.Type == "image" {
		fileName = "files/images/"
	} else {
		log.Printf("Unsupported type")
		ErrorResponseStatus(context, http.StatusBadRequest, "Unsupported type")
		return
	}

	// save the file in image
	err = context.SaveUploadedFile(file, fileName+file.Filename)
	if err != nil {
		log.Printf("Saving File Error: ", err)
		ErrorResponseStatus(context, http.StatusInternalServerError, "Can't save the file")
		return
	}

	// save the resource in the data base
	err = database.InsertNewResource(database.Ctx, reqResource, file, fileName)

	if err != nil {
		log.Println(err)
		ErrorResponseStatus(context, http.StatusInternalServerError, "Failed to store the resource")
		return
	}

	SuccessResponseStatus(context, http.StatusOK, "successfully stored the resource")
}

// // Get the secret Information
// // @Summary      analyze the secret
// // @Description  This routes generate new secret have the user's data
// // @Tags         token
// // @Accept       json
// // @Accept       xml
// // @Produce      json
// // @Produce      xml
// // @Param        token  path      string  true  "get the secret info"
// // @Success      200  {object}    models.ResponseData
// // @Failure      400  {object}	  models.ErrorModel
// // @Failure      500  {object}	  models.ErrorModel
// // @Router       /get/{token} [post]
// func GetToken(context *gin.Context) {
// 	tokenString := context.Param("token")

// 	claims := jwt.MapClaims{}
// 	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(os.Getenv("SECRET_TOKEN")), nil
// 	})

// 	if err != nil {
// 		ErrorResponseStatus(context, http.StatusBadRequest, "invalid token")
// 		return
// 	}

// 	info := token.Claims.(*jwt.MapClaims)
// 	id := (*info)["id"]
// 	secretData := (*info)["data"]
// 	objectId, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", id))
// 	if err != nil {
// 		ErrorResponseStatus(context, http.StatusBadRequest, "invalid id")
// 		return
// 	}
// 	// get the token object from database
// 	var object models.Secret
// 	err = database.SearchCollection.
// 		FindOne(database.Ctx, bson.D{{Key: "_id", Value: objectId}}).
// 		Decode(&object)

// 	if err != nil {
// 		ErrorResponseStatus(context, http.StatusBadRequest, "invalid token")
// 		return
// 	}

// 	if object.Views <= 0 {
// 		// delete expired object
// 		_, err := database.SearchCollection.DeleteOne(database.Ctx, bson.D{{Key: "_id", Value: objectId}})
// 		if err != nil {
// 			ErrorResponseStatus(context, http.StatusInternalServerError, "database error")
// 			return
// 		}
// 		ErrorResponseStatus(context, http.StatusBadRequest, "No views available")
// 		return
// 	}
// 	// update the number of views
// 	filter := bson.D{{Key: "_id", Value: objectId}}
// 	update := bson.D{{Key: "$set", Value: bson.D{{Key: "views", Value: object.Views - 1}}}}
// 	_, err = database.SearchCollection.UpdateOne(database.Ctx, filter, update)

// 	if err != nil {
// 		ErrorResponseStatus(context, http.StatusInternalServerError, "database error")
// 		return
// 	}
// 	SuccessResponseStatusForToken(context, http.StatusOK, secretData, object)
// }
