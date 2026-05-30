include .env
export

proto-user:
	protoc \
	  -I proto/user \
	  --go_out=./services/user/pkg/pb --go_opt=paths=source_relative \
	  --go-grpc_out=./services/user/pkg/pb --go-grpc_opt=paths=source_relative \
	  proto/user/*.proto

proto-upload:
	protoc \
	  -I proto/upload \
	  --go_out=./services/upload/pkg/pb --go_opt=paths=source_relative \
	  --go-grpc_out=./services/upload/pkg/pb --go-grpc_opt=paths=source_relative \
	  proto/upload/*.proto

proto-userpb-upload:
	protoc \
	  -I proto/user \
	  --go_out=./services/upload/pkg/userpb --go_opt=paths=source_relative \
	  --go-grpc_out=./services/upload/pkg/userpb --go-grpc_opt=paths=source_relative \
	  proto/user/*.proto

proto-uploadpb-transcode:
	protoc \
	  -I proto/upload \
	  --go_out=./services/transcode/pkg/uploadpb --go_opt=paths=source_relative \
	  --go-grpc_out=./services/transcode/pkg/uploadpb --go-grpc_opt=paths=source_relative \
	  proto/upload/*.proto

proto-all:
	make proto-user
	make proto-upload
	make proto-userpb-upload
	make proto-uploadpb-transcode

docker-up:
	docker compose --env-file .env up --build

docker-down:
	docker compose down