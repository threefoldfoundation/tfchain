all: install

daemonpkgs = ./cmd/tfchaind
clientpkgs = ./cmd/tfchainc
testpkgs = ./pkg/types
pkgs = $(daemonpkgs) $(clientpkgs) ./pkg/config $(testpkgs)

version = $(shell git describe | cut -d '-' -f 1)
commit = $(shell git rev-parse --short HEAD)
ifeq ($(commit), $(shell git rev-list -n 1 $(version) | cut -c1-7))
fullversion = $(version)
else
fullversion = $(version)-$(commit)
endif

dockerVersion = $(shell git describe | cut -d '-' -f 1| cut -d 'v' -f 2)
dockerVersionEdge = edge

configpkg = github.com/threefoldfoundation/tfchain/pkg/config
ldflagsversion = -X $(configpkg).rawVersion=$(fullversion)

stdoutput = $(GOPATH)/bin
daemonbin = $(stdoutput)/tfchaind
clientbin = $(stdoutput)/tfchainc

install:
	go build -race -tags='debug profile' -ldflags '$(ldflagsversion)' -o $(daemonbin) $(daemonpkgs)
	go build -race -tags='debug profile' -ldflags '$(ldflagsversion)' -o $(clientbin) $(clientpkgs)

install-std:
	go build -ldflags '$(ldflagsversion)' -o $(daemonbin) $(daemonpkgs)
	go build -ldflags '$(ldflagsversion)' -o $(clientbin) $(clientpkgs)

test:
	go test -race -v -tags='debug testing' -timeout=60s $(testpkgs)

# xc builds and packages release binaries
# for all windows, linux and mac, 64-bit only,
# using the standard Golang toolchain.
xc:
	docker build -t tfchainbuilder -f DockerBuilder .
	docker run --rm -v $(shell pwd):/go/src/github.com/threefoldfoundation/tfchain tfchainbuilder

docker-minimal: xc
	docker build -t tfchain/tfchain:$(dockerVersion) -f DockerfileMinimal --build-arg binaries_location=release/tfchain-$(version)-linux-amd64/cmd .

# Release images builds and packages release binaries, and uses the linux based binary to create a minimal docker
release-images: get_hub_jwt docker-minimal
	docker push tfchain/tfchain:$(dockerVersion)
	# also create a latest
	docker tag tfchain/tfchain:$(dockerVersion) tfchain/tfchain
	docker push tfchain/tfchain:latest
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X POST --data "image=tfchain/tfchain:$(dockerVersion)" "https://hub.gig.tech/api/flist/me/docker"
	# symlink the latest flist
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X GET "https://hub.gig.tech/api/flist/me/tfchain-tfchain-$(dockerVersion).flist/link/tfchain-tfchain-latest.flist"
	# Merge the flist with ubuntu and nmap flist, so we have a tty file etc...
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X POST --data "[\"gig-official-apps/ubuntu1604.flist\", \"tfchain/tfchain-tfchain-$(dockerVersion).flist\", \"gig-official-apps/nmap.flist\"]" "https://hub.gig.tech/api/flist/me/merge/ubuntu-16.04-tfchain-$(dockerVersion).flist"
	# And also link in a latest
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X GET "https://hub.gig.tech/api/flist/me/ubuntu-16.04-tfchain-$(dockerVersion).flist/link/ubuntu-16.04-tfchain-latest.flist"

xc-edge:
	docker build -t tfchainbuilderedge -f DockerBuilderEdge .
	docker run --rm -v $(shell pwd):/go/src/github.com/threefoldfoundation/tfchain tfchainbuilderedge

docker-minimal-edge: xc-edge
	docker build -t tfchain/tfchain:$(dockerVersionEdge) -f DockerfileMinimal --build-arg binaries_location=release/tfchain-$(dockerVersionEdge)-linux-amd64/cmd .

# Release images builds and packages release binaries, and uses the linux based binary to create a minimal docker
release-images-edge: get_hub_jwt docker-minimal-edge
	docker push tfchain/tfchain:$(dockerVersionEdge)
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X POST --data "image=tfchain/tfchain:$(dockerVersionEdge)" "https://hub.gig.tech/api/flist/me/docker"
	# Merge the flist with ubuntu and nmap flist, so we have a tty file etc...
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X POST --data "[\"gig-official-apps/ubuntu1604.flist\", \"tfchain/tfchain-tfchain-$(dockerVersionEdge).flist\", \"gig-official-apps/nmap.flist\"]" "https://hub.gig.tech/api/flist/me/merge/ubuntu-16.04-tfchain-$(dockerVersionEdge).flist"

explorer: release-dir
	tar -C ./frontend -czvf release/explorer-$(dockerVersion).tar.gz explorer

explorer-edge: release-dir
	tar -C ./frontend -czvf release/explorer-$(dockerVersionEdge).tar.gz explorer

release-dir:
	[ -d release ] || mkdir release

release-explorer: get_hub_jwt explorer
	# Upload explorer
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -F file=@./release/explorer-$(dockerVersion).tar.gz "https://hub.gig.tech/api/flist/me/upload"
	# Symlink latest
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X GET "https://hub.gig.tech/api/flist/me/explorer-$(dockerVersion).flist/link/explorer-latest.flist"
	# Merge with caddy
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X POST --data "[\"gig-official-apps/caddy.flist\", \"tfchain/explorer-$(dockerVersion).flist\"]" "https://hub.gig.tech/api/flist/me/merge/caddy-explorer-$(dockerVersion).flist"
	# And also link in a latest
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X GET "https://hub.gig.tech/api/flist/me/caddy-explorer-$(dockerVersion).flist/link/caddy-explorer-latest.flist"

release-explorer-edge: get_hub_jwt explorer-edge
	# Upload explorer
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -F file=@./release/explorer-$(dockerVersionEdge).tar.gz "https://hub.gig.tech/api/flist/me/upload"
	# Merge with caddy
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X POST --data "[\"gig-official-apps/caddy.flist\", \"tfchain/explorer-$(dockerVersionEdge).flist\"]" "https://hub.gig.tech/api/flist/me/merge/caddy-explorer-$(dockerVersionEdge).flist"

get_hub_jwt: check-HUB_APP_ID check-HUB_APP_SECRET
	$(eval HUB_JWT = $(shell curl -X POST "https://itsyou.online/v1/oauth/access_token?response_type=id_token&grant_type=client_credentials&client_id=$(HUB_APP_ID)&client_secret=$(HUB_APP_SECRET)&scope=user:memberof:tfchain"))

check-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Required env var $* not present"; \
		exit 1; \
	fi

ineffassign:
	ineffassign $(pkgs)

.PHONY: all install xc release-images get_hub_jwt check-% ineffassign explorer release-explorer