#!/bin/sh
# Uploads, instantiates, and executes a wasm contract + queries

export KEY="eve1" # eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn
export KEYALGO="secp256k1"
export EVED_CHAIN_ID="eve-t1"
export EVED_KEYRING_BACKEND="os"
export EVED_NODE="http://localhost:26657"
export EVED_COMMAND_ARGS="--gas-prices="0.025ueve" --gas 5000000 -y --from $KEY"

export KEY_ADDR="eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn"
echo "decorate bright ozone fork gallery riot bus exhaust worth way bone indoor calm squirrel merry zero scheme cotton until shop any excess stage laundry" | eved keys add $KEY --keyring-backend $EVED_KEYRING_BACKEND --algo $KEYALGO --recover


# For some reason eve returns all values in base64 & not human readable. Need to fix that
getCodeID() {
    TX_HASH=$1
    # eved q tx $TX_HASH --output json | jq -r '.events[] | select(.type=="store_code").attributes[0].value' | base64 --decode
    eved q tx $TX_HASH --output json | jq -r '.logs[].events[] | select(.type=="store_code").attributes[0].value'
}

echo "STORING CW721 TO CHAIN"
TX721=$(eved tx wasm store cw721_base.wasm -y --broadcast-mode block --output json $EVED_COMMAND_ARGS | jq -r '.txhash')
# CODE_ID_721=$(eved query tx $TX721 --output json | jq -r '.logs[0].events[-1].attributes[0].value')
CODE_ID_721=`getCodeID $TX721`
echo "CW721 WAS STORED, WITH CODE ID $CODE_ID_721"
NFT721_TX_UPLOAD=$(eved tx wasm instantiate "$CODE_ID_721" '{"name": "eved-721","symbol": "ctest","minter": "eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn"}' --label "base_cw_721" $EVED_COMMAND_ARGS --output json --broadcast-mode block -y --admin $KEY_ADDR | jq -r '.txhash')
echo "INSTANCE instantiate'd"
ADDR721=$(eved query tx $NFT721_TX_UPLOAD --output json | jq -r '.logs[0].events[0].attributes[0].value') && echo "ADDR 721: $ADDR721"
# eve1wkwy0xh89ksdgj9hr347dyd2dw7zesmtrue6kfzyml4vdtz6e5wsj3vejy

function mintToken() {
    CONTRACT_ADDR=$1
    TOKEN_ID=$2
    OWNER=$3
    TOKEN_URI=$4

    export EXECUTED_MINT_JSON=`printf '{"mint":{"token_id":"%s","owner":"%s","token_uri":"%s"}}' $TOKEN_ID $OWNER $TOKEN_URI`
    TXMINT=$(eved tx wasm execute "$CONTRACT_ADDR" "$EXECUTED_MINT_JSON" --from $KEY --yes --output json --broadcast-mode block | jq -r '.txhash') && echo $TXMINT
}

# mint a token with id 1 to the contract
mintToken $ADDR721 1 "eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn" "https://ipfs.io/ipfs/QmNLijobERK4VhSDZdKjt5SrezdRM6k813qcSHd68f3Mqg"

# query the token data directly
echo $(eved q wasm contract-state smart "$ADDR721" '{"all_nft_info":{"token_id":"1"}}' --output json) | jq -r '.data.info.token_uri'

# query all token_ids in a given contract
eved query wasm contract-state smart $ADDR721 '{"tokens":{"owner":"eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn","start_after":"0","limit":50}}'