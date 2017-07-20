
.PHONY: tserver clean all

all: tserver

tserver :
	cd tserver && \
	go build

clean :
	rm tserver/tserver
