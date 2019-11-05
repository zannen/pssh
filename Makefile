# Makefile

test:
	@go test ./...

bench:
	@go test -run=X -bench=. -benchmem ./...

