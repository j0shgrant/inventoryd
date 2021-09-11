default: lint install

install: install-inventoryctl install-inventoryd

install-inventoryctl:
	go install ./inventoryctl

install-inventoryd:
	go install ./inventoryd

lint:
	golangci-lint run
