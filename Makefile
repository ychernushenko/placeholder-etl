.PHONY: build test clean run

test:
	$(MAKE) clean
	docker-compose run --rm placeholder-etl go test -coverprofile=coverage.out ./...

build:
	docker-compose build --no-cache

clean:
	docker-compose down --volumes --rmi all
	rm -rf ./data
	rm -rf ./logs

run:
	$(MAKE) clean
	$(MAKE) build
	docker-compose up --build
