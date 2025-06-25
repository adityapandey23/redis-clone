run: build
	@./bin/redis-clone --listenAddr :5000

build: clean
	@go build -o bin/redis-clone .

clean:
	@go clean -testcache

