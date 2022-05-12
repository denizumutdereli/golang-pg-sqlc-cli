dockerdbs:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
pgstart:
	docker start postgres
createdb:
	docker exec -it postgres createdb --encoding=UTF8 --username=root --owner=root mservice
dropdb:
	docker exec -it postgres dropdb mservice
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@127.0.0.1:5432/mservice?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@127.0.0.1:5432/mservice?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./db/... -timeout 10s
.PHONY: dockerdbs pgstart createdb dropdb migrateup migratedown sqlc test