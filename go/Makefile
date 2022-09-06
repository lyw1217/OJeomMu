default: build

buildTime=$(shell date -u "+%Y%m%d%I%M")
revVersion=$(shell git rev-parse --short HEAD)
build: 
	go mod download && go mod verify;
	go build -v -o ../../bin/ojeommu.$(revVersion).$(buildTime) .