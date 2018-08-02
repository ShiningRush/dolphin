install:
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...

gen:install
	go-bindata --nocompress -pkg dolphinui -o dist/dolphinui/datafile.go  thirdparty/dashboard/...

build: gen
	go build
	cd ./task && go build
	cd ./syncer && go build
	cd ./dashserver && go build

test: build
	cd ./task && go test
	cd ./syncer && go test

clean:
	rm -rf dist/dolphinui

.PHONY: install build test clean