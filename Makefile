# Makefile

install:
	@go install .

test:
	@go test ./...

bench:
	@go test -run=X -bench=. -benchmem ./...
