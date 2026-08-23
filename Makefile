build:
   go build -o bin/snap

test:
   go test ./...

lint:
   golangci-lint run

vet:
   go vet ./...

tidy:
   go mod tidy