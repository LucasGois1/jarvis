createmigration:
	migrate create -ext=sql -dir=sql/migrations -seq init

migrate:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/jarvis" -verbose up
migratedown:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/jarvis" -verbose drop
grpc:
	protoc --go_out=. --go-grpc_out=. proto/chat.proto --experimental_allow_proto3_optional
.PHONY: 
	migrate createmigration migratedown grpc