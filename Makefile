# ==============================================================================
# Help

.PHONY: help
## help: shows this help message
help:
	@ echo "Usage: make [target]\n"
	@ sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==============================================================================
# Tests

.PHONY: test
## test: run unit tests
test:
	@ go test -v -race ./... -count=1

.PHONY: coverage
## coverage: run unit tests and generate coverage report in html format
coverage:
	@ packages=$$(go list ./... | grep -v "cmd" | grep -v "validate" | grep -v "gui"); \
	if [ -z "$$packages" ]; then \
		echo "No valid Go packages found"; \
		exit 1; \
	fi; \
	go test -race -coverpkg=$$(echo $$packages | tr ' ' ',') -coverprofile=coverage.out $$packages && go tool cover -html=coverage.out

# ==============================================================================
# Building

.PHONY: build-gui-app
## build-gui-app: build the GUI application
build-gui-app:
	@ go build -o cmd/gui/MacOSDMGCreator cmd/gui/main.go

.PHONY: build-sample-app
## build-sample-app: build the sample application
build-sample-app:
	@ go build -o dmg/integration/sampleapp/SampleApp dmg/integration/sampleapp/cmd/main.go

# ==============================================================================
# Quality Checks

.PHONY: vet
## vet: runs Go vet to analyze code for potential issues
vet:
	@ echo "Running go vet..."
	@ go vet ./...

.PHONY: govulncheck
## govulncheck: runs Go vulnerability check
govulncheck:
	@ go install golang.org/x/vuln/cmd/govulncheck@latest
	@ echo "Running go vuln check..."
	@ govulncheck ./...

# ==============================================================================
# App execution

.PHONY: run-gui
## run-gui: run the GUI application
run-gui:
	@ go run cmd/gui/main.go
