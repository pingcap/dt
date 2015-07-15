all: build

build: build-cmd build-instance

build-cmd: 
	go build -o bin/cmd ./cmd

build-instance:
	go build -o bin/instance ./sample

clean:
	@rm -rf bin
	@rm sample/*.out
	@rm sample/*.log
