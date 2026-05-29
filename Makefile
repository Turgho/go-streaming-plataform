include .env
export

proto-user:
	protoc \
	  -I proto/user \
	  --go_out=./services/user/pkg/pb --go_opt=paths=source_relative \
	  --go-grpc_out=./services/user/pkg/pb --go-grpc_opt=paths=source_relative \
	  proto/user/*.proto

docker-up:
	docker compose --env-file .env up --build

docker-down:
	docker compose down