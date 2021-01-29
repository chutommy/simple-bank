.PHONY: postgres createdb dropdb newmigration migrateup migratedown sqlc test

postgres:
	docker run -p 5432:5432 --env POSTGRES_PASSWORD=simplebankpassword --name postgres12 -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it postgres12 dropdb --username=postgres simple_bank

# make create-migration FILENAME=[migration name]
newmigration:
	docker run --rm -v $(PWD)/db/migration:/migrations --network host migrate/migrate create -ext sql -dir /migrations -seq $(FILENAME)

migrateup:
	docker run --rm -v $(PWD)/db/migration:/migrations --network host migrate/migrate -path /migrations -database "postgresql://postgres:simplebankpassword@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	docker run --rm -v $(PWD)/db/migration:/migrations --network host migrate/migrate -path /migrations -database "postgresql://postgres:simplebankpassword@localhost:5432/simple_bank?sslmode=disable" -verbose drop -f

sqlc:
	docker run --rm -v $(PWD):/src -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...
