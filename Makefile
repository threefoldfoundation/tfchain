all: install

pkgs = ./cmd/tfchainc ./cmd/tfchaind
version = $(shell git describe | cut -d '-' -f 1)
versionTag = $(shell git describe | cut -d '-' -f 1| cut -d 'v' -f 2)

install:
	go install -race -tags='debug profile' $(pkgs)

# xc builds and packages release binaries
# for all windows, linux and mac, 64-bit only,
# using the standard Golang toolchain.
xc:
	docker build -t tfchainbuilder -f DockerBuilder .
	docker run --rm -v $(shell pwd):/go/src/github.com/threefoldfoundation/tfchain tfchainbuilder

docker-minimal: xc
	docker build -t tfchain/tfchain:$(versionTag) -f DockerfileMinimal --build-arg binaries_location=release/tfchain-$(version)-linux-amd64/cmd .

# Release images builds and packages release binaries, and uses the linux based binary to create a minimal docker
release-images: get_hub_jwt docker-minimal
	docker push tfchain/tfchain:$(versionTag)
	# also create a latest
	docker tag tfchain/tfchain:$(versionTag) tfchain/tfchain
	docker push tfchain/tfchain:latest
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X POST --data "image=tfchain/tfchain:$(versionTag)" "https://hub.gig.tech/api/flist/me/docker"
	# symlink the latest flist
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X GET "https://hub.gig.tech/api/flist/me/tfchain-tfchain-$(versionTag).flist/link/tfchain-tfchain-latest.flist"
	# Merge the flist with ubuntu flist, so we have a tty file etc...
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X POST --data "sources=[gig-official-apps/ubuntu1604.flist, tfchain/tfchain-tfchain-$(versionTag).flist]" "https://hub.gig.tech/api/flist/me/merge/ubuntu-16.04-tfchain-$(versionTag).flist"
	# And also link in a latest
	curl -b "active-user=tfchain; caddyoauth=$(HUB_JWT)" -X GET "https://hub.gig.tech/api/flist/me/ubuntu-16.04-tfchain-$(versionTag).flist/link/ubuntu-16.04-tfchain-latest.flist"

get_hub_jwt: check-HUB_APP_ID check-HUB_APP_SECRET
	$(eval HUB_JWT = $(shell curl -X POST "https://itsyou.online/v1/oauth/access_token?response_type=id_token&grant_type=client_credentials&client_id=$(HUB_APP_ID)&client_secret=$(HUB_APP_SECRET)&scope=user:memberof:tfchain"))

check-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Required env var $* not present"; \
		exit 1; \
	fi

.PHONY: all install xc release-images get_hub_jwt check-%