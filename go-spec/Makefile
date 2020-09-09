test:
	go test -cover -race ./...

generate_proto:
	echo "make sure you have protoc-go-gen v1.3.5 ONLY!"
	protoc -I=./ --go_out=./src/core ./src/core/*.proto
