.PHONY: proto-gen
proto-gen:
	protoc --proto_path=laptop/proto --go_out=laptop --go_opt=paths=source_relative --go-grpc_out=laptop --go-grpc_opt=paths=source_relative laptop/proto/*.proto
