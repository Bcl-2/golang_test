dbUp:
	docker run --name test-sql -e MYSQL_ROOT_PASSWORD=qwerty123 -p 3063:3063 -d mysql:latest

dbDown:
	docker stop test-sql
	docker rm test-sql

build: main.go
	go build main.go
	./main


