#!/bin/bash

KEY="test"
CHAINID="eve-testnet"
KEYRING="test"
MONIKER="local-testnet"
KEYALGO="secp256k1"
LOGLEVEL="info"
VALIDATOR="validator"
TESTING_ACCOUNT="vesting_account"


echo >&1 "installing eve"
rm -rf $HOME/.eved
make install

eve config keyring-backend $KEYRING
eve config chain-id $CHAINID

# determine if user wants to recorver or create new

eved keys add $VALIDATOR --keyring-backend $KEYRING
MY_VALIDATOR_ADDRESS=$(eved keys show $VALIDATOR -a --keyring-backend $KEYRING)

eved keys add $TESTING_ACCOUNT --keyring-backend $KEYRING
MY_VESTING_ACCOUNT=$(eved keys show $TESTING_ACCOUNT -a --keyring-backend $KEYRING)

echo >&1 "\n"

# init chain
eved init $MONIKER --chain-id $CHAINID

# Change parameter token denominations to ubaby
cat $HOME/.eved/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="ueve"' > $HOME/.eved/config/tmp_genesis.json && mv $HOME/.eved/config/tmp_genesis.json $HOME/.eved/config/genesis.json
cat $HOME/.eved/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="ueve"' > $HOME/.eved/config/tmp_genesis.json && mv $HOME/.eved/config/tmp_genesis.json $HOME/.eved/config/genesis.json
cat $HOME/.eved/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="ueve"' > $HOME/.eved/config/tmp_genesis.json && mv $HOME/.eved/config/tmp_genesis.json $HOME/.eved/config/genesis.json
cat $HOME/.eved/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["amount"]="100000"' > $HOME/.eved/config/tmp_genesis.json && mv $HOME/.eved/config/tmp_genesis.json $HOME/.eved/config/genesis.json
cat $HOME/.eved/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["max_deposit_period"]="200s"' > $HOME/.eved/config/tmp_genesis.json && mv $HOME/.eved/config/tmp_genesis.json $HOME/.eved/config/genesis.json
cat $HOME/.eved/config/genesis.json | jq '.app_state["gov"]["voting_params"]["voting_period"]="250s"' > $HOME/.eved/config/tmp_genesis.json && mv $HOME/.eved/config/tmp_genesis.json $HOME/.eved/config/genesis.json
cat $HOME/.eved/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="ueve"' > $HOME/.eved/config/tmp_genesis.json && mv $HOME/.eved/config/tmp_genesis.json $HOME/.eved/config/genesis.json

# Allocate genesis accounts (cosmos formatted addresses)
eved add-genesis-account $VALIDATOR 1000000000ueve --keyring-backend $KEYRING

# Sign genesis transaction
eved gentx $VALIDATOR 700000000ueve --keyring-backend $KEYRING --chain-id $CHAINID

# Collect genesis tx
eved collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
eved validate-genesis

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
eved start --minimum-gas-prices=0.0001ueve