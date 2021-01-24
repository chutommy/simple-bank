.PHONY: postgres createdb dropdb newmigration migrateup migratedown

postgres:
	docker run -p 5433:5432 --env POSTGRES_PASSWORD=simplebankpassword --name postgres12 postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it postgres12 dropdb --username=postgres simple_bank

# make create-migration FILE=[migration name]
newmigration:
	migrate create -ext sql -dir db/migration -seq $(FILE)

migrateup:
	migrate -path db/migration -database "postgresql://postgres:simplebankpassword@localhost:5433/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:simplebankpassword@localhost:5433/simple_bank?sslmode=disable" -verbose down
