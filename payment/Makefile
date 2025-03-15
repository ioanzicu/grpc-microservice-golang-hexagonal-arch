.PHONY: test
test:
	go clean -testcache
	go test -race -v -cover ./...

.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: upgrade
upgrade:
	go get -u ./...
	go mod tidy

.PHONY: docker-build
docker-build:
	docker build -t orders .

.PHONY: upgrade-dependencies
upgrade-dependencies:
    go get -u ./... 
	go mod tidy