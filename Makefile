gen:
	go generate ./...
dev: 
	make gen
	SECRET=secret go run .
test:
	go test -v ./...
test.cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
lint:
	find . -type f -name "*.templ" -exec templ fmt "{}" \;
	gofumpt -d -w .
	golangci-lint run -v
