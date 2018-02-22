#!/bin/bash
set -e

full_version=$(git describe | cut -d '-' -f 1,3)
version=$(echo "$full_version" | cut -d '-' -f 1)

for os in darwin linux windows; do
	echo Packaging ${os}...
	# create workspace
	folder="release/tfchain-${version}-${os}-amd64"
	rm -rf "$folder"
	mkdir -p "$folder"
	# compile and sign binaries
	for pkg in cmd/tfchainc cmd/tfchaind; do
		bin=$pkg
		if [ "$os" == "windows" ]; then
			bin=${pkg}.exe
		fi
		GOOS=${os} go build -a -tags 'netgo' \
			-ldflags="-s -w" \
			-o "${folder}/${bin}" "./${pkg}"

	done
	# add other artifacts
	cp -r LICENSE README.md "$folder"
	# zip
	(
		zip -rq "release/tfchain-${version}-${os}-amd64.zip" \
			"release/tfchain-${version}-${os}-amd64"
	)
done