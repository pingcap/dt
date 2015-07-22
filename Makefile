all: build

build: build-cmd build-instance

build-cmd: 
	go build -o bin/agent ./cmd
	go build -o bin/ctrl ./cmd

build-instance:
	go build -o bin/probe ./sample/probe
	go build -o bin/instance ./sample/instance

clean:
	@rm -rf bin
	@rm sample/*.out
	@rm sample/*.log
	@rm -rf sample/dt
