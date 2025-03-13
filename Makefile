.PHONY: build test clean run

test:
	docker-compose run --rm placeholder-etl go test ./...

build:
	docker-compose build --no-cache

clean:
	docker-compose down --volumes --rmi all
	rm -rf ./data

run:
	$(MAKE) clean
	$(MAKE) build
	docker-compose up --build
