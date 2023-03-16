package controllers

import (
	"log"
	"net/http"
	"os"

	"searchEngine/database"
	"searchEngine/models"

	"github.com/gin-gonic/gin"
)

// ErrorResponseStatus function using return response with error status
func ErrorResponseStatus(context *gin.Context, status int, message string) {
	context.JSON(status, gin.H{"message": message})
}

// SuccessResponseStatus function using return a response with Ok status
func SuccessResponseStatus(context *gin.Context, status int, data string) {
	context.JSON(status, gin.H{"message": data})
}

type LoggerLevels struct {
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
	DebugLogger *log.Logger
}

type ResourceService struct {
	Database database.Resource
	Logger   LoggerLevels
}

func NewResourceService(database database.Resource) *ResourceService {
	return &ResourceService{
		Database: database,
		Logger: LoggerLevels{
			ErrorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
			InfoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
			DebugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		},
	}
}

// CreateResource Function used to Create new resource
// @Summary      Create a new resource
// @Description  This route uses to create a new resource it accepts form data with your file form data contains Title as string, Type as string (only supprt this values["HTML Pages", "image"]), and Tags as array of string
// @Accept       mpfd
// @Produce      json
// @Param 		 file  formData file true "File"
// @Param        data  formData models.ReqResource true "Create new resource"
// @Success      200  {object}  models.CreateResourceResponse
// @Failure      400  {object}	models.CreateResourceResponse
// @Failure      500  {object}	models.CreateResourceResponse
// @Router       /createNewResource [post]
func (r *ResourceService) CreateResource(context *gin.Context) {
	r.Logger.InfoLogger.Println("controllers,", "resource.go,", "CreateResource() Func")

	var reqResource models.ReqResource
	if context.PostForm("title") != "" {
		reqResource.Title = context.PostForm("title")
	} else {
		r.Logger.ErrorLogger.Println("title is required")
		ErrorResponseStatus(context, http.StatusBadRequest, "title is required")
		return
	}

	if context.PostForm("type") != "" {
		reqResource.Type = context.PostForm("type")
	} else {
		r.Logger.ErrorLogger.Println("type is required")
		ErrorResponseStatus(context, http.StatusBadRequest, "type is required")
		return
	}

	if len(context.PostFormArray("tags")) != 0 {
		reqResource.Tags = context.PostFormArray("tags")
	} else {
		r.Logger.ErrorLogger.Println("tags is required")
		ErrorResponseStatus(context, http.StatusBadRequest, "tags is required")
		return
	}
	file, err := context.FormFile("file")
	if err != nil {
		r.Logger.ErrorLogger.Println("File Error: ", err)
		ErrorResponseStatus(context, http.StatusBadRequest, "Can't upload the file")
		return
	}

	var storedPath string
	var realStoredPath string

	if reqResource.Type == "HTML Page" {
		storedPath = "../files/HTML/"
		realStoredPath = "searchEngine/files/HTML/"

	} else if reqResource.Type == "image" {
		storedPath = "../files/images/"
		realStoredPath = "searchEngine/files/images/"
	} else {
		r.Logger.ErrorLogger.Println("Unsupported type, the type should be either HTML Page or image")
		ErrorResponseStatus(context, http.StatusBadRequest, "Unsupported type, the type should be either HTML Page or image")
		return
	}

	/*
		**********************************************************
		*********************** Note *****************************
		I assumed that the frontend part should send the file to the backend,
		then store it.

		In this example, I store the files and images in the docker image.
		It's not good because the size of the file will be large but I did it to show you how to store it,
		and I don't have another choice.

		The best way is to use platforms to store files and images such as Amazon S3,
		then you can render these files and image via url.
	*/
	err = context.SaveUploadedFile(file, storedPath+file.Filename)
	if err != nil {
		r.Logger.ErrorLogger.Println("Saving File Error: ", err)
		ErrorResponseStatus(context, http.StatusInternalServerError, "Can't save the file")
		return
	}

	// save the resource in the database
	err = r.Database.InsertNewResource(database.Ctx, reqResource, file.Filename, realStoredPath)

	if err != nil {
		r.Logger.ErrorLogger.Println("Failed to store the resource: ", err)
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
		r.Logger.ErrorLogger.Println("missing keyword: ", err)
		ErrorResponseStatus(context, http.StatusBadRequest, "missing keyword")
		return
	}

	result, err := r.Database.GetResources(database.Ctx, searchEngineRequest.Keyword)
	if err != nil {
		r.Logger.ErrorLogger.Println("failed to get data from the database: ", err)
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
		r.Logger.ErrorLogger.Println("missing id")
		ErrorResponseStatus(context, http.StatusBadRequest, "missing id")
		return
	}

	result, err := r.Database.GetResource(database.Ctx, id)
	if err != nil {
		r.Logger.ErrorLogger.Println("failed to get data from the database: ", err)
		ErrorResponseStatus(context, http.StatusInternalServerError, "failed to get data from the database")
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": result})

}
