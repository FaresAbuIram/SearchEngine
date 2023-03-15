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

type ResourceService struct {
	Database database.Resource
}

func NewResourceService(database database.Resource) *ResourceService {
	return &ResourceService{
		Database: database,
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
func (r *ResourceService) CreateResource(context *gin.Context) {
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
	err = r.Database.InsertNewResource(database.Ctx, reqResource, file.Filename, fileName)

	if err != nil {
		log.Println(err)
		ErrorResponseStatus(context, http.StatusInternalServerError, "Failed to store the resource")
		return
	}

	SuccessResponseStatus(context, http.StatusOK, "successfully stored the resource")
}

// Search for a keyword
// @Summary      Search for a keyword
// @Description  This route uses to Search for a keyword in a tags
// @Accept       json
// @Produce      json
// @Param        body  body models.SearchEngineRequest true "Search for a resource"
// @Success      200  {object}  []models.SearchEngineResult
// @Failure      400  {object}	models.CreateResourceResponse
// @Failure      500  {object}	models.CreateResourceResponse
// @Router       /search [post]
func (r *ResourceService) Search(context *gin.Context) {
	var searchEngineRequest models.SearchEngineRequest
	if err := context.BindJSON(&searchEngineRequest); err != nil {
		log.Printf("missing keyword: ", err)
		ErrorResponseStatus(context, http.StatusBadRequest, "missing keyword")
		return
	}

	result, err := r.Database.GetResources(database.Ctx, searchEngineRequest.Keyword)
	if err != nil {
		log.Printf("failed to get data from the database: ", err)
		ErrorResponseStatus(context, http.StatusInternalServerError, "failed to get data from the database")
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": result})
}

// Get resource information
// @Summary      Get resource information
// @Description  This route uses to Get resource information in order to visualise the resource
// @Accept       json
// @Produce      json
// @Param        id  path string true "resource id"
// @Success      200  {object}  models.Resource
// @Failure      400  {object}	models.CreateResourceResponse
// @Failure      500  {object}	models.CreateResourceResponse
// @Router       /resource/{id} [get]
func (r *ResourceService) GetResource(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		log.Printf("missing id")
		ErrorResponseStatus(context, http.StatusBadRequest, "missing id")
		return
	}

	result, err := r.Database.GetResource(database.Ctx, id)
	if err != nil {
		log.Printf("failed to get data from the database: ", err)
		ErrorResponseStatus(context, http.StatusInternalServerError, "failed to get data from the database")
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": result})

}
