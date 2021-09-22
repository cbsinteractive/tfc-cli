.PHONY: fmt staticcheck test

fmt:
	@go fmt ./...

staticcheck: fmt
	@staticcheck ./...

test:
	@go test ./...
