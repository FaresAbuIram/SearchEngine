mocks:
	@mockery@2.14.0 --name="Resource" --dir="./database" --output="./controllers/mocks"

tests:
	go test controllers/resource_test.go 
