#migrate -path db/migration -database "postgresql://root:secret@192.168.99.100:5432/transaction_bp?sslmode=disable" -verbose up
postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it postgres createdb --username=root --owner=root transaction_bp
dropdb:
	docker exec -it dropdb transaction_bp
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@192.168.99.100:5432/transaction_bp?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@192.168.99.100:5432/transaction_bp?sslmode=disable" -verbose down
sqlc:
	sqlc generate
.PHONY: postgres createdb dropdb migrateup migratedown sqlc