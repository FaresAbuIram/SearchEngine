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
	router := gin.Default()
	resourceMocks := new(ResourceMocks.Resource)
	resourceService := controllers.NewResourceService(resourceMocks)
	ReqResource := models.ReqResource{
		Title: "title",
		Type:  "HTML Pages",
		Tags:  []string{"tag1", "tag2"},
	}

	resourceMocks.On("InsertNewResource", mock.AnythingOfType("*context.emptyCtx"), ReqResource, "example.html", "files/HTML/").Return(nil)

	router.POST("/createNewResource", resourceService.CreateResource)
	file, err := os.Open("example.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	formData := map[string]string{
		"title": "title",
		"type":  "HTML Pages",
	}
	tags := []string{"tag1", "tag2"}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	filePart, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	_, err = io.Copy(filePart, file)
	if err != nil {
		fmt.Println(err)
		// panic(err)
	}

	for key, value := range formData {
		_ = writer.WriteField(key, value)
	}

	for _, tag := range tags {
		writer.WriteField("tags", tag)
	}

	req, err := http.NewRequest("POST", "/createNewResource", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	err = writer.Close()
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", rr.Code)
	}

	expected := `{"message":"successfully stored the resource"}`
	if rr.Body.String() != expected {
        t.Errorf("Expected %q but got %q", expected, rr.Body.String())
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

	resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
        t.Errorf("Expected status code %d but got %d", http.StatusOK, resp.Code)
    }

	expected := `{"data":[{"ID":"1","Title":"Fares"}]}`
    if resp.Body.String() != expected {
        t.Errorf("Expected %q but got %q", expected, resp.Body.String())
    }
}

func TestGetResource(t *testing.T) {
	router := gin.Default()
	resourceMocks := new(ResourceMocks.Resource)
	resourceService := controllers.NewResourceService(resourceMocks)

	response := models.Resource{
		Title: "Fares",
		Type: "HTML Pages",
		Tags: []string{"tag1", "tag2"},
		Path: "files/HTML/example.html",
	}

	resourceMocks.On("GetResource", mock.AnythingOfType("*context.emptyCtx"), "1").Return(response, nil)
	router.POST("/resource/:id", resourceService.GetResource)

	req, err := http.NewRequest("POST", "/resource/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
        t.Errorf("Expected status code %d but got %d", http.StatusOK, resp.Code)
    }
	
	expected := `{"data":{"Title":"Fares","Type":"HTML Pages","Tags":["tag1","tag2"],"Path":"files/HTML/example.html"}}`
    if resp.Body.String() != expected {
        t.Errorf("Expected %q but got %q", expected, resp.Body.String())
    }
}
