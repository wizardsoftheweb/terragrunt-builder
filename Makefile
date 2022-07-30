# Aliases for executables
GO ?= go
RM ?= rm

test:
	$(GO) test -v ./... -cover -race -coverprofile=.coverage.out

coverage: test
	$(GO) tool cover -func=.coverage.out

coverage-report: coverage
	$(GO) tool cover -html=.coverage.out

clean:
	rm -rf .coverage.out
