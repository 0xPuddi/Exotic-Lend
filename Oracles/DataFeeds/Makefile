run: build
	@./bin/DataFeedExec

build:
	@rm -f ./bin/DataFeedExec && go build -C ./src -o ../bin/DataFeedExec

tidy:
	@cd src && go mod tidy

postgredev:
	@docker run --name pgsql-dev --rm -e POSTGRES_PASSWORD=dev -p 7777:7777 -d postgres

test:
	@cd src && go test ./...

test-cover:
	@rm -f ./bin/coverage.out && \
	cd src && \
	go test ./... -coverprofile=../bin/coverage.out && \
	echo "\n" && \
	go tool cover -func=../bin/coverage.out

test-cover-view:
	@rm -f ./bin/coverage.out ./bin/coverage.html && cd src && \
	go test ./... -coverprofile=../bin/coverage.out && \
	go tool cover -html=../bin/coverage.out -o ../bin/coverage.html && \
	open ../bin/coverage.html