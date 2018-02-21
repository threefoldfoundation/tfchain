all: install

pkgs = ./cmd/tfchainc ./cmd/tfchaind

install:
	go install -race -tags='debug profile' $(pkgs)