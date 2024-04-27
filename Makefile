PHONY: generate-structs
generate-structs:
	mkdir "./pkg/user_v1"
	protoc --go_out=pkg/user_v1 --go_opt=paths=source_relative \
	api/user_v1/service.proto


PHONY: generate

generate:
	mkdir pkg\user_v1
	protoc --go_out=pkg/user_v1 --go_opt=paths=import \
			--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=import \
			api/user_v1/service.proto
	move pkg\user_v1\github.com\Kerupie\DiplomaV2\pkg\user_v1\* pkg\user_v1
	rmdir /s /q pkg\user_v1\github.com