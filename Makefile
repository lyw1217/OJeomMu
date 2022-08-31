default: build

build: 
	go mod download && go mod verify
	go build -v -o ./ojeommu .