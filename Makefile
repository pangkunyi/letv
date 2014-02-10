export GOPATH=$(shell pwd)

install:
	@go install letv
run:
	@./bin/letv
test:
	@./bin/letv -id 2256635 -res 1080p
