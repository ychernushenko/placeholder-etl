.PHONY: build up down logs test

test:
	docker-compose run --rm placeholder-etl go test ./...

run:
	docker-compose up --build
