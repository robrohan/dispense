
HASH = $(shell git log --pretty=format:'%h' -n 1)

# List all targets in thie file
list:
	@echo ""
	@echo "~ Dispense ~"
	@echo ""
	@grep -B 1 '^[^#[:space:]].*:' Makefile
	@echo ""

# List all libraries
install:
	go mod tidy

# Test the application
test:
	@go test ./...

# Run the app locally
run:
	@mkdir -p public
	@go run cmd/dispense/main.go

clean:
	rm -rf build

# Build the app to distribute
build: clean
	@mkdir -p build
	@go build -o build/dispense \
		-ldflags "-X main.build=$(HASH)" \
		cmd/dispense/main.go
