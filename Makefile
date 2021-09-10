default: build

build: build-inventoryctl build-inventoryd

build-inventoryctl:
	go build -o "${GOPATH}/bin" ./inventoryctl

build-inventoryd:
	go build -o "${GOPATH}/bin" ./inventoryd
