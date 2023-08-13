initdb :
	docker run --name postgres15 -v /data/db-data/postgres15:/var/lib/postgresql/data -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:15-alpine 

postgres:
	docker exec -it postgres15 psql

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root chat

dropdb:
	docker exec -it postgres15 dropdb go-chat

.PHONY:	initDB postgres createdb dropdb



build:
	@go build -o bin/server 

run: build 
	@./bin/server
