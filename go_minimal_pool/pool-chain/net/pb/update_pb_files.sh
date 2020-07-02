rm ./pool-chain/net/pb/*.pb.go


protoc -I=./pool-chain/net/pb/ --go_out=plugins=grpc:./ ./pool-chain/net/pb/*.proto
