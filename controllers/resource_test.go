package controllers_test

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"searchEngine/controllers"
	ResourceMocks "searchEngine/controllers/mocks"
	"searchEngine/models"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

func TestCreateResource(t *testing.T) {
	// prepare required data
	router := gin.Default()
	resourceMocks := new(ResourceMocks.Resource)
	resourceService := controllers.NewResourceService(resourceMocks)

	// fake data 
	ReqResource := models.ReqResource{
		Title: "title",
		Type:  "HTML Page",
		Tags:  []string{"tag1", "tag2"},
	}

	// mock the InsertNewResource function with fake data
	resourceMocks.On("InsertNewResource", mock.AnythingOfType("*context.emptyCtx"), ReqResource, "example.html", "searchEngine/files/HTML/").Return(nil)

	// prepare our request
	router.POST("/createNewResource", resourceService.CreateResource)

	// create file form
	file, err := os.Open("example.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// create form data
	formData := map[string]string{
		"title": "title",
		"type":  "HTML Page",
	}
	tags := []string{"tag1", "tag2"}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	filePart, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// write our fake data to the body
	_, err = io.Copy(filePart, file)
	if err != nil {
		panic(err)
	}
	for key, value := range formData {
		_ = writer.WriteField(key, value)
	}
	for _, tag := range tags {
		writer.WriteField("tags", tag)
	}

	// prapare the request before send it
	req, err := http.NewRequest("POST", "/createNewResource", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	err = writer.Close()
	if err != nil {
		panic(err)
	}


	response := httptest.NewRecorder()
	// send the request
	router.ServeHTTP(response, req)


	// Test cases
	if response.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", response.Code)
	}
	expected := `{"message":"successfully stored the resource"}`
	if response.Body.String() != expected {
        t.Errorf("Expected %q but got %q", expected, response.Body.String())
    }
}

// the request will missing the title 
func TestCreateResourceNegative(t *testing.T) {
	// prepare required data
	router := gin.Default()
	resourceMocks := new(ResourceMocks.Resource)
	resourceService := controllers.NewResourceService(resourceMocks)

	// fake data 
	ReqResource := models.ReqResource{
		Type:  "HTML Page",
		Tags:  []string{"tag1", "tag2"},
	}

	// mock the InsertNewResource function with fake data
	resourceMocks.On("InsertNewResource", mock.AnythingOfType("*context.emptyCtx"), ReqResource, "example.html", "searchEngine/files/HTML/").Return(nil)

	// prepare our request
	router.POST("/createNewResource", resourceService.CreateResource)

	// create file form
	file, err := os.Open("example.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// create form data
	formData := map[string]string{
		"type":  "HTML Page",
	}
	tags := []string{"tag1", "tag2"}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	filePart, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// write our fake data to the body
	_, err = io.Copy(filePart, file)
	if err != nil {
		panic(err)
	}
	for key, value := range formData {
		_ = writer.WriteField(key, value)
	}
	for _, tag := range tags {
		writer.WriteField("tags", tag)
	}

	// prapare the request before send it
	req, err := http.NewRequest("POST", "/createNewResource", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	err = writer.Close()
	if err != nil {
		panic(err)
	}


	response := httptest.NewRecorder()
	// send the request
	router.ServeHTTP(response, req)


	// Test cases
	if response.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest; got %v", response.Code)
	}
	expected := `{"message":"title is required"}`
	if response.Body.String() != expected {
        t.Errorf("Expected %q but got %q", expected, response.Body.String())
    }
}

func TestSearch(t *testing.T) {
	router := gin.Default()
	resourceMocks := new(ResourceMocks.Resource)
	resourceService := controllers.NewResourceService(resourceMocks)

	response := models.SearchEngineResult{
		ID:    "1",
		Title: "Fares",
	}

	resourceMocks.On("GetResources", mock.AnythingOfType("*context.emptyCtx"), "string").Return([]models.SearchEngineResult{response}, nil)

	router.POST("/search", resourceService.Search)

	reqBody := `{"keyword": "string"}`
	req, err := http.NewRequest("POST", "/search", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	result := httptest.NewRecorder()
    router.ServeHTTP(result, req)
	if result.Code != http.StatusOK {
        t.Errorf("Expected status code %d but got %d", http.StatusOK, result.Code)
    }

	expected := `{"data":[{"ID":"1","Title":"Fares"}]}`
    if result.Body.String() != expected {
        t.Errorf("Expected %q but got %q", expected, result.Body.String())
    }
}

func TestGetResource(t *testing.T) {
	router := gin.Default()
	resourceMocks := new(ResourceMocks.Resource)
	resourceService := controllers.NewResourceService(resourceMocks)

	response := models.Resource{
		Title: "Fares",
		Type: "HTML Page",
		Tags: []string{"tag1", "tag2"},
		Path: "searchEngine/files/HTML/example.html",
	}

	resourceMocks.On("GetResource", mock.AnythingOfType("*context.emptyCtx"), "1").Return(response, nil)
	router.POST("/resource/:id", resourceService.GetResource)

	req, err := http.NewRequest("POST", "/resource/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	result := httptest.NewRecorder()
    router.ServeHTTP(result, req)
	if result.Code != http.StatusOK {
        t.Errorf("Expected status code %d but got %d", http.StatusOK, result.Code)
    }
	
	expected := `{"data":{"Title":"Fares","Type":"HTML Page","Tags":["tag1","tag2"],"Path":"searchEngine/files/HTML/example.html"}}`
    if result.Body.String() != expected {
        t.Errorf("Expected %q but got %q", expected, result.Body.String())
    }
}
