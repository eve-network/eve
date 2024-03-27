#!/bin/bash

BINARY=build/eved

CHAIN_ID=local-eve

CONTRACT=cw20_base
HOME=mytestnet
DENOM="ueve"
KEY="test0"

KEYRING="test"
TEST0_ADDRESS=$($BINARY keys show $KEY -a --keyring-backend $KEYRING --home $HOME)

if ! $BINARY tx wasm store scripts/wasm/$CONTRACT.wasm --from test0 --keyring-backend $KEYRING --chain-id $CHAIN_ID --gas auto --gas-adjustment 1.3 --home $HOME -y; then
    echo "Error uploading contract"
    exit 1
fi

# $BINARY query wasm list-contract-by-code $PROPOSAL
# $BINARY query wasm list-contracts-by-creator $TEST0_ADDRESS
$BINARY query wasm list-code