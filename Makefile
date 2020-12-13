GOFLAGS = CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GOTEST_PACKAGES = $(shell go list ./... | egrep -v '(pkg|cmd)')

gotest:
	go test -race -v -cover -coverprofile coverage.out $(GOTEST_PACKAGES)

lint:
	golangci-lint run -v

coclient:
	docker exec -it dbhook_cockroachdb_1 ./cockroach sql --insecure 
	
generate:
	go generate ./...
