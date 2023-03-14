package models

type ReqResource struct {
	Title string   `bson:"title" binding:"required"`
	Type  string   `bson:"type" binding:"required"`
	Tags  []string `bson:"tags" binding:"required"`
}

type Resource struct {
	Title string   `bson:"title" binding:"required"`
	Type  string   `bson:"type" binding:"required"`
	Tags  []string `bson:"tags" binding:"required"`
	Path  string   `bson:"path" binding:"required"`
}

type CreateResourceResponse struct {
	Message string `bson:"message"`
}

type SearchEngineRequest struct {
	keyword    string `bson:"keyword"`
}

type SearchEngineResult struct {
	ID    string `bson:"_id"`
	Title string `bson:"title"`
}
