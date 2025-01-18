run: build
	@./bin/goredis --listenAddr :4000

build: 
	@go build -o bin/goredis .