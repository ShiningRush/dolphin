install:
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...

gen:
	go-bindata --nocompress -pkg dolphinui -o dist/dolphinui/datafile.go  thirdparty/dashboard/...

build: gen
	go build
	cd ./task && go build
	cd ./syncer && go build
	cd ./dashserver && go build

clean:
	rm -rf dist/dolphinui

.PHONY: install build clean