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

OS_LIST=(linux darwin)

ARCHIVE=false
if [ "$1" = "archive" ]; then
	ARCHIVE=true
	shift # remove element from arguments
fi

# Overide the file names to edge version, keep full version at the git commit since
# that is the expected format
if [ "$1" = "edge" ]; then
	version="edge"
	# if more params defined, use them as OS_LIST
	if [[ "$#" -gt 1 ]]; then
		OS_LIST=("${@:2}")
	fi
elif [[ "$#" -ge 1 ]]; then
	OS_LIST=("${@:1}")
fi

# ensure xgo is installed
go get -u github.com/karalabe/xgo

# Create temp work space
tmpfolder="release/tfchain-xc.tmp"
rm -rf "$tmpfolder"
mkdir -p "$tmpfolder"

# Compile targets list
TARGETS=""
for os in "${OS_LIST[@]}"; do
	TARGETS+="${os}/amd64,"
done
TARGETS="${TARGETS%?}"

# Compile binaries
for pkg in ./cmd/tfchainc ./cmd/tfchaind ./cmd/bridged; do
	xgo --go 1.11.x --targets="${TARGETS}" \
		-ldflags="-X ${package}/pkg/config.rawVersion=${full_version} -s -w" \
		-out "$(basename $pkg)-$version" \
		-dest "$tmpfolder" \
		$pkg
done

if [ "$ARCHIVE" = false ] ; then
    exit 0 # finished already
fi

# Create archives
for os in "${OS_LIST[@]}"; do
	folder="release/tfchain-${version}-${os}-amd64"
	rm -rf "$folder"
	mkdir -p "$folder/cmd"

	# copy binaries
	for binary in $(ls "$tmpfolder" | grep "$os" | grep "amd64"); do
		cp "${tmpfolder}/${binary}" "$folder/cmd/$(echo "$binary" | cut -d- -f1)"
	done

	# copy other artifacts
	cp -r doc LICENSE README.md "$folder"

	# go into the release directory
	pushd release &> /dev/null
	# zip
	ziparchive="tfchain-${version}-${os}-amd64.zip"
	rm -f "$ziparchive"
	(
		zip -rq "$ziparchive" "tfchain-${version}-${os}-amd64"
	)
	# leave the release directory
	popd &> /dev/null

	# clean up workspace dir
	rm -rf "$folder"
done

# all archives are ready clean up temp directory
rm -rf "$tmpfolder"
