all: install

daemonpkgs = ./cmd/tfchaind
clientpkgs = ./cmd/tfchainc
bridgepkgs = ./cmd/bridged
faucetpkgs = ./frontend/tftfaucet
testpkgs = ./pkg/types ./pkg/persist ./pkg/eth ./pkg/cli
pkgs = $(daemonpkgs) $(clientpkgs) ./pkg/config $(testpkgs)

version = $(shell git describe --abbrev=0)
commit = $(shell git rev-parse --short HEAD)
ifeq ($(commit), $(shell git rev-list -n 1 $(version) | cut -c1-7))
	fullversion = $(version)
	fullversionpath = \/releases\/tag\/$(version)
else
	fullversion = $(version)-$(commit)
	fullversionpath = \/tree\/$(commit)
endif

dockerVersion = $(shell git describe --abbrev=0 | cut -d 'v' -f 2)
dockerVersionEdge = edge

configpkg = github.com/threefoldfoundation/tfchain/pkg/config
ldflagsversion = -X $(configpkg).rawVersion=$(fullversion)

stdoutput = $(GOPATH)/bin
daemonbin = $(stdoutput)/tfchaind
clientbin = $(stdoutput)/tfchainc
bridgebin = $(stdoutput)/bridged

install:
	go build -race -tags='debug profile' -ldflags '$(ldflagsversion)' -o $(daemonbin) $(daemonpkgs)
	go build -race -tags='debug profile' -ldflags '$(ldflagsversion)' -o $(clientbin) $(clientpkgs)
	go build -race -tags='debug profile' -ldflags '$(ldflagsversion)' -o $(bridgebin) $(bridgepkgs)

install-std:
	go build -ldflags '$(ldflagsversion) -s -w' -o $(daemonbin) $(daemonpkgs)
	go build -ldflags '$(ldflagsversion) -s -w' -o $(clientbin) $(clientpkgs)

update:
	git pull && git submodule update --recursive --remote

test:
	go test -race -v -tags='debug testing' -timeout=60s $(testpkgs)

test-coverage:
	gocoverutil -coverprofile cover.out test \
		-short -race -v -tags='debug testing' -timeout=60s -covermode=atomic \
		$(testpkgs)

test-coverage-web: test-coverage
	go tool cover -html=cover.out

# xc builds and packages release binaries
# for all windows, linux and mac, 64-bit only,
# using the standard Golang toolchain.
xc:
	bash release.sh

docker-minimal: xc
	docker build -t tfchain/tfchain:$(dockerVersion) -f DockerfileMinimal --build-arg binaries_location=release/tfchain-$(version)-linux-amd64/cmd .

# Release images builds and packages release binaries, and uses the linux based binary to create a minimal docker
release-images: get_hub_jwt docker-minimal
	docker push tfchain/tfchain:$(dockerVersion)
	# also create a latest
	docker tag tfchain/tfchain:$(dockerVersion) tfchain/tfchain
	docker push tfchain/tfchain:latest
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X POST --data "image=tfchain/tfchain:$(dockerVersion)" "https://hub.grid.tf/api/flist/me/docker"
	# symlink the latest flist
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X GET "https://hub.grid.tf/api/flist/me/tfchain-tfchain-$(dockerVersion).flist/link/tfchain-tfchain.flist"
	# Merge the flist with ubuntu and nmap flist, so we have a tty file etc...
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X POST --data "[\"tf-official-apps/ubuntu1604.flist\", \"tfchain/tfchain-tfchain-$(dockerVersion).flist\", \"tf-official-apps/nmap.flist\"]" "https://hub.grid.tf/api/flist/me/merge/ubuntu-16.04-tfchain-$(dockerVersion).flist"
	# And also link in a latest
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X GET "https://hub.grid.tf/api/flist/me/ubuntu-16.04-tfchain-$(dockerVersion).flist/link/ubuntu-16.04-tfchain.flist"

xc-edge:
	bash release.sh edge

docker-minimal-edge: xc-edge
	docker build -t tfchain/tfchain:$(dockerVersionEdge) -f DockerfileMinimal --build-arg binaries_location=release/tfchain-$(dockerVersionEdge)-linux-amd64/cmd .

# Release images builds and packages release binaries, and uses the linux based binary to create a minimal docker
release-images-edge: get_hub_jwt docker-minimal-edge
	docker push tfchain/tfchain:$(dockerVersionEdge)
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X POST --data "image=tfchain/tfchain:$(dockerVersionEdge)" "https://hub.grid.tf/api/flist/me/docker"
	# Merge the flist with ubuntu and nmap flist, so we have a tty file etc...
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X POST --data "[\"tf-official-apps/ubuntu1604.flist\", \"tfchain/tfchain-tfchain-$(dockerVersionEdge).flist\", \"tf-official-apps/nmap.flist\"]" "https://hub.grid.tf/api/flist/me/merge/ubuntu-16.04-tfchain-$(dockerVersionEdge).flist"

explorer: release-dir embed-explorer-version
	tar -C $(TEMPDIR)/frontend -czvf release/explorer-$(dockerVersion).tar.gz explorer
	rm -r $(TEMPDIR)

explorer-edge: release-dir embed-explorer-version
	tar -C $(TEMPDIR)/frontend -czvf release/explorer-$(dockerVersionEdge).tar.gz explorer
	rm -r $(TEMPDIR)

embed-explorer-version:
	$(eval TEMPDIR = $(shell mktemp -d))
	cp -r ./frontend $(TEMPDIR)
	sed -i '' 's/version=0/version=$(fullversion)/g' $(TEMPDIR)/frontend/explorer/public/*.html
	sed -i '' 's/version=null/version=\"$(fullversion)\"/g' $(TEMPDIR)/frontend/explorer/public/js/footer.js
	sed -i '' 's/versionpath=null/versionpath=\"$(fullversionpath)\"/g' $(TEMPDIR)/frontend/explorer/public/js/footer.js

release-dir:
	[ -d release ] || mkdir release

release-explorer: get_hub_jwt explorer
	# Upload explorer
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -F file=@./release/explorer-$(dockerVersion).tar.gz "https://hub.grid.tf/api/flist/me/upload"
	# Symlink latest
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X GET "https://hub.grid.tf/api/flist/me/explorer-$(dockerVersion).flist/link/explorer.flist"
	# Merge with caddy
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X POST --data "[\"tf-official-apps/caddy.flist\", \"tfchain/explorer-$(dockerVersion).flist\"]" "https://hub.grid.tf/api/flist/me/merge/caddy-explorer-$(dockerVersion).flist"
	# And also link in a latest
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X GET "https://hub.grid.tf/api/flist/me/caddy-explorer-$(dockerVersion).flist/link/caddy-explorer.flist"

faucet: release-dir
	GOOS=linux go build -o ./release/faucet $(faucetpkgs)

release-faucet: faucet get_hub_jwt
	tar -C ./release -cvzf release/faucet.tar.gz faucet
	# Upload to hub
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -F file=@./release/faucet.tar.gz "https://hub.grid.tf/api/flist/me/upload"
	# Merge with caddy
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X POST --data "[\"tf-official-apps/caddy.flist\", \"tfchain/faucet.flist\"]" "https://hub.grid.tf/api/flist/me/merge/caddy-faucet.flist"

release-explorer-edge: get_hub_jwt explorer-edge
	# Upload explorer
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -F file=@./release/explorer-$(dockerVersionEdge).tar.gz "https://hub.grid.tf/api/flist/me/upload"
	# Merge with caddy
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X POST --data "[\"tf-official-apps/caddy.flist\", \"tfchain/explorer-$(dockerVersionEdge).flist\"]" "https://hub.grid.tf/api/flist/me/merge/caddy-explorer-$(dockerVersionEdge).flist"

get_hub_jwt: check-HUB_APP_ID check-HUB_APP_SECRET
	$(eval HUB_JWT = $(shell curl -X POST "https://itsyou.online/v1/oauth/access_token?response_type=id_token&grant_type=client_credentials&client_id=$(HUB_APP_ID)&client_secret=$(HUB_APP_SECRET)&scope=user:memberof:tfchain"))

check-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Required env var $* not present"; \
		exit 1; \
	fi

ineffassign:
	ineffassign $(pkgs)

.PHONY: all install xc release-images get_hub_jwt check-% ineffassign explorer release-explorer faucet