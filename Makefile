run: build
	@./bin/redis-clone --listenAddr :5002

build: clean
	@go build -o bin/redis-clone .

clean:
	@go clean -testcache

