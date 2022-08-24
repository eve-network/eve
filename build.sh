#/bin/sh

# VERSION=$(echo $(git describe --tags) | sed 's/^v//')
VERSION="0.0.1"
COMMIT="$(git log -1 --format='%H')"

# put this in makefile later
go install -ldflags "-w -s -X github.com/cosmos/cosmos-sdk/types.DBBackend=pebbledb -X github.com/cosmos/cosmos-sdk/version.Name=eved -X github.com/cosmos/cosmos-sdk/version.Version=$VERSION -X github.com/cosmos/cosmos-sdk/version.Commit=$COMMIT" -tags pebbledb ./...