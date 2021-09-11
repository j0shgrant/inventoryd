default: dependencies install

install: install-inventoryctl install-inventoryd

install-inventoryctl:
	go install ./inventoryctl

install-inventoryd:
	go install ./inventoryd

dependencies:
	go mod download
