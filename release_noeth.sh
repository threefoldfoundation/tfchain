#!/bin/bash
set -e

package="github.com/threefoldfoundation/tfchain"

version=$(git describe --abbrev=0)
commit="$(git rev-parse --short HEAD)"
if [ "$commit" == "$(git rev-list -n 1 $version | cut -c1-7)" ]
then
	full_version="$version"
else
	full_version="${version}-${commit}"
fi

# Overide the file names to edge version, keep full version at the git commit since
# that is the expected format
if [ "$1" = "edge" ]; then
	version="edge"
fi

echo "building version ${version}"

for os in darwin linux; do
	echo Packaging ${os}...
	# create workspace
	folder="release/tfchain-noeth-${version}-${os}-amd64"
	rm -rf "$folder"
	mkdir -p "$folder"
	# compile and sign binaries
	for pkg in cmd/tfchainc cmd/tfchaind; do
		GOOS=${os} go build -a \
			-ldflags="-X ${package}/pkg/config.rawVersion=${full_version} -s -w" \
			-tags="noeth" -o "${folder}/${pkg}" "./${pkg}"

	done

	# add other artifacts
	cp -r doc LICENSE README.md "$folder"
	# go into the release directory
	pushd release &> /dev/null
	# zip
	(
		zip -rq "tfchain-noeth-${version}-${os}-amd64.zip" \
			"tfchain-noeth-${version}-${os}-amd64"
	)
	# leave the release directory
	popd &> /dev/null

	# clean up workspace dir
	rm -rf "$folder"
done
