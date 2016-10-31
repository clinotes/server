run:
		ENV=local go run *.go

test-data:
		@go test github.com/clinotes/server/data -v

test: test-data
