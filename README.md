# Search engine and simplified annotations

## Description
This search engine is developed using web annotation.When users enter specific words or phrases in a search engine, it automatically fetches the most relevant resources that contain those keywords. Web annotation makes it possible. Web annotation helps to make an application user-friendly. Thanks to web annotation, users can add, modify, and remove information from Web resources without altering the resource itself.

This project uses web annotation on pages and images. When the user enters words,
names, or phrases in the system, it will fetch the information and pictures having the same annotation. Then the system displays a list of results that contain the image or content matching to the user input. For this search engine, you need to use an effective algorithm to generate a query result page/search result records based on users’ queries.


## How to run it
### Using Docker
- Clone the repo. 
- Run this command `make up` (make sure you are in the home dir)
- Open the swagger docs, to test the apps `http://localhost:8080/swagger/index.html`

### Locally 
- Clone the repo. 
- Setup the MongoDB Locally 
- Run `go run main.go` command
- Open the swagger docs, to test the apps `http://localhost:8080/swagger/index.html`
