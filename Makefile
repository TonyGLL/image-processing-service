initschema:
	migrate create -ext sql -dir sql/migrations -seq init_schema

postgres:
	docker compose up -d --build db

stop-postgres:
	docker stop postgres_image_processing

createdb:
	docker exec -it postgres_image_processing createdb --username=root --owner=root image_processing

dropdb:
	docker exec -it postgres_image_processing dropdb image_processing

migrateup:
	migrate -path sql/migrations -database "postgresql://root:secret@localhost:5432/image_processing?sslmode=disable" -verbose up

migratedown:
	migrate -path sql/migrations -database "postgresql://root:secret@localhost:5432/image_processing?sslmode=disable" -verbose down

build:
	@echo "Building..."

	@CONFIG_FILE=local.env go build -o main cmd/app/main.go

	@echo "Build successfully!!!"

start: build
	CONFIG_FILE=local.env air

.PHONY: initschema postgres stop-postgres createdb dropdb migrateup migratedown