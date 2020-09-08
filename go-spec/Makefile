test:
	rm -f ./src/mocks/mock_types.go
	mockgen -source=./src/core/types.go -destination=./src/mocks/mock_types.go -package=mocks
	go test -cover -race ./...
