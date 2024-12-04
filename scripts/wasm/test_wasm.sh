#!/bin/bash

# run chain locally, get contract if not exists, upload contract, instantiate contract, query contract

# check if have screen
if ! [ -x "$(command -v screen)" ]; then
    echo 'Error: screen is not installed. Please install using apt-get install screen (Linux) or brew install screen(MacOS)' >&2
    exit 1
fi

# Run chain
echo "Starting chain..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    screen -L -dmS node1 bash scripts/run-node.sh
else
    screen -L -Logfile $HOME/log-screen.txt -dmS node1 bash scripts/run-node.sh
fi
sleep 30
# check if already have wasm, if not get contract
if [ -f "./scripts/wasm/cw20_base.wasm" ]; then
    echo "Wasm file already exists for cw20_base"
else
    echo "Getting contract..."
    if ! bash ./scripts/wasm/get_contract.sh; then
        echo "Error getting contract"
        exit 1
    fi
    sleep 30
fi

# upload contract and check if error
echo "Uploading contract..."
if ! bash ./scripts/wasm/upload_code.sh; then
    echo "Error uploading contract"
    exit 1
fi

sleep 30

# instantiate contract and check if error
echo "Instantiating contract..."
if ! bash ./scripts/wasm/instantiate_cw20.sh; then
    echo "Error instantiating and querying contract..."
    exit 1
fi

echo "Success executing basic wasm tests!"

# kill chain
pkill -9 limed



