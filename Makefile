gen-grpc:
	protoc -I=api/proto --go_out=. api/proto/getInfo.proto
	protoc -I=api/proto --go-grpc_out=. api/proto/getInfo.proto

gen-gateway:
	protoc -I=api/proto \
		--grpc-gateway_out ./api/proto \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt grpc_api_configuration=api/proto/getInfo.yaml \
		--openapiv2_out ./api/swagger \
		--openapiv2_opt logtostderr=true \
		--openapiv2_opt generate_unbound_methods=true \
		api/proto/getInfo.proto

gen-all: gen-grpc gen-gateway gen-swagger

up:
	docker compose up
down:
	docker-compose down

