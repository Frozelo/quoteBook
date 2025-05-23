APP_NAME = quote-book
PKG = ./...
DOCKER_IMAGE = quote-book

test:
	go test -v ./internal/store ./internal/handlers

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker run --rm -p 8080:8080 $(DOCKER_IMAGE)

docker-self-check:
	docker run --rm -e SELF_CHECK=1 $(DOCKER_IMAGE)

check: docker-build docker-self-check

.PHONY: test docker-build docker-run docker-self-check check
