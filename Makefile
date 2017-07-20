
.PHONY: tserver clean tserver_test test all

all: tserver

tserver :
	cd tserver && \
	go build

test : tserver_test

tserver_test : tserver
	cd tserver && \
	go test

clean :
	rm tserver/tserver
