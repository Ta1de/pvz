run:
	go run cmd/main.go

PKGS := $(shell go list ./... | grep -vE '/test')
COVERPKG := $(shell go list ./... | grep -vE '/test' | paste -sd, -)

cover_func:
	@echo "Running tests with coverage..."
	@rm -f coverage.out
	@go test -coverprofile=coverage.out -coverpkg=$(COVERPKG) $(PKGS)
	@go tool cover -func=coverage.out

cover_html:
	@echo "Running tests with coverage..."
	@rm -f coverage.out
	@go test -coverprofile=coverage.out -coverpkg=$(COVERPKG) $(PKGS)
	@go tool cover -html=coverage.out

inter_test:
	go test test/integration_test.go

docker_run:
	docker compose up -d