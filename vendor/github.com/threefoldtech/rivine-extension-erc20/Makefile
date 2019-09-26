testpkgs = ./api/bridge ./types
allpkgs = ./api ./api/bridge ./client ./daemon ./examples/erc20_exchange_wallet ./http ./types .

test: ineffassign
	go test -race -v -tags='debug testing' -timeout=60s $(testpkgs)

ineffassign:
	ineffassign ./api ./api/bridge ./client ./daemon ./examples/erc20_exchange_wallet ./http ./types .

lint:
	goimports -w $(allpkgs)
	gofmt -s -w $(allpkgs)

.PHONY: test ineffassign lint
