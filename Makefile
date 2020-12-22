postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migrateup:
	migrate-go -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate-go -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
migrateup-mysql:
	migrate-go -path db/migration -database "mysql://root:jianxin@tcp(localhost:3306)/windy" -verbose up
migratedown-mysql:
	migrate-go -path db/migration -database "mysql://root:jianxin@tcp(localhost:3306)/windy" -verbose down
migrateinit:
	migrate-go create -ext sql -dir db/migration -seq init-schema
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown migrateinit sqlc server