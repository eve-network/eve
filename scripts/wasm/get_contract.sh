#!/bin/bash

VERSION="v0.16.0"
CONTRACTS="cw20_base"

for CONTRACT in $CONTRACTS; do
  # check if already exists wasm
  if [ -f ./scripts/wasm/$CONTRACT.wasm ]; then
    echo "Wasm file already exists for $CONTRACT"
    continue
  fi
  curl -s -L -o ./scripts/wasm/$CONTRACT.wasm https://github.com/CosmWasm/cw-plus/releases/download/$VERSION/$CONTRACT.wasm
done