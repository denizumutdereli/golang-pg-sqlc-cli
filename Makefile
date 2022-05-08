postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it postgres createdb --encoding=UTF8 --username=root --owner=root mservice
dropdb:
	docker exec -it postgres dropdb mservice
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@192.168.99.100:5432/mservice?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@192.168.99.100:5432/mservice?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./db/...
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test