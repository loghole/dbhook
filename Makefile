GOFLAGS = CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GOTEST_PACKAGES = $(shell go list ./... | egrep -v '(pkg|cmd)')

.PHONY: test
test:
	go test -race -v -cover -coverprofile coverage.out $(GOTEST_PACKAGES)

.PHONY: lint
lint:
	golangci-lint run ./... -c .golangci.yaml -v 

.PHONY: coclient
coclient:
	docker exec -it dbhook_cockroachdb_1 ./cockroach sql --insecure 
	
.PHONY: generate
generate:
	go generate ./...

.PHONY: imports
imports:
	goimports --local "github.com/loghole/dbhook" -w -d $$(find . -type f -name '*.go')
