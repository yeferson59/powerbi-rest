.PHONY: run db down stop

run:
	@go run main.go

db:
	@docker-compose up -d db

down:
	@docker-compose down

stop:
	@docker-compose stop
