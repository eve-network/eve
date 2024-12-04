#!/bin/bash

# this is cw20 code from install_contracts.sh, true if this is the first proposal
CODE=1
CHAIN_ID=local-lime
BINARY=build/limed
HOME=mytestnet
$BINARY keys add demo --keyring-backend test --home $HOME
VAL_KEY=test0
VAL=$($BINARY keys show -a $VAL_KEY --keyring-backend test --home $HOME)
DEMO=$($BINARY keys show -a demo --keyring-backend test --home $HOME)


# Init params in json
INIT=$(cat <<EOF
{
  "name": "My first token",
  "symbol": "FRST",
  "decimals": 6,
  "initial_balances": [{
    "address": "$VAL",
    "amount": "123456789000"
  }]
}
EOF
)

if ! $BINARY tx wasm instantiate $CODE "$INIT" --label "First Coin" --no-admin --from $VAL_KEY --keyring-backend test --chain-id $CHAIN_ID -y --gas auto --gas-adjustment 1.3 --home $HOME; then
    echo "Error instantiating contract"
    exit 1
fi

# wait the chain to process the tx
echo "Waiting for tx to be processed..."
sleep 5
# query the contract address, include this in case of multiple contract instantiation, default getting the first contract
CONTRACT_POSITION=$1
if [ -z "$CONTRACT_POSITION" ]; then
    CONTRACT_POSITION=0
fi
CONTRACT=$($BINARY q wasm list-contracts-by-creator "$VAL" --output json | jq -r ".contract_addresses[$CONTRACT_POSITION]" )
if [ -z "$CONTRACT" ]; then
    echo "No contract found"
    exit 1
fi

QUERY=$(cat <<EOF
{ "balance": { "address": "$VAL" }}
EOF
)
QUERY_DEMO=$(cat <<EOF
{ "balance": { "address": "$DEMO" }}
EOF
)

# check initial balance
echo "Validator Balance:"
$BINARY query wasm contract-state smart $CONTRACT "$QUERY"
echo "Demo Balance:"
$BINARY query wasm contract-state smart $CONTRACT "$QUERY_DEMO"

# send some tokens
TRANSFER=$(cat <<EOF
{
  "transfer": {
    "recipient": "$DEMO",
    "amount": "987654321"
  }
}
EOF
)
# $BINARY tx wasm execute $CONTRACT "$TRANSFER" --from $VAL_KEY --keyring-backend test \
#     --chain-id $CHAIN_ID -y --gas auto --gas-adjustment 1.3 --home $HOME

if ! $BINARY tx wasm execute $CONTRACT "$TRANSFER" --from $VAL_KEY --keyring-backend test --chain-id $CHAIN_ID -y --gas auto --gas-adjustment 1.3 --home $HOME; then
    echo "Error executing contract"
    exit 1
fi
# wait the chain to process the tx
echo "Waiting for tx to be processed..."
sleep 5

# check final balance
echo "Validator Balance:"
$BINARY query wasm contract-state smart $CONTRACT "$QUERY"
echo "Demo Balance:"
$BINARY query wasm contract-state smart $CONTRACT "$QUERY_DEMO"