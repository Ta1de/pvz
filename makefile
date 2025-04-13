run:
	go run cmd/main.go

PKGS := $(shell go list ./... | grep -vE '/test')
COVERPKG := $(shell go list ./... | grep -vE '/test' | paste -sd, -)

cover_test:
	@echo "Running tests with coverage..."
	@rm -f coverage.out
	@go test -coverprofile=coverage.out -coverpkg=$(COVERPKG) $(PKGS)
	@go tool cover -func=coverage.out

cover_html:
	go test -coverprofile=coverage.out -coverpkg=./... ./... | go tool cover -html=coverage.out


