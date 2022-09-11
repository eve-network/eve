# call on build to make a new folder which we can mount too via compose.
# this way we can build each node to share the latests gentx
# most commands taken from ./test_node.sh

CHAINID="eve-d"  # eve-docker
KEYALGO="secp256k1"
KEYRING="test"  # export EVE_KEYRING="TEST"
LOGLEVEL="info"
KEY="eve1"      # eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn
MONIKER1="dockereve1"
KEY2="eve2"     # eve1j4rtuq6zm5mmw9xcjmm7gymlj39tvwnt9h4sm2
MONIKER2="dockereve2"

# File Location
SCRIPT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
TESTNETS_SHARED_DIR=$SCRIPT_DIR/.testnets
# check that SCRIPT_DIR is longer than 7, if so delete the .testnets folder for a fresh start
if [ ${#SCRIPT_DIR} -gt 7 ]; then
  echo "$SCRIPT_DIR"
  rm -rf $SCRIPT_DIR/.testnets
fi

V1=$SCRIPT_DIR/.testnets/v1
V2=$SCRIPT_DIR/.testnets/v2

# init chain from key 1, giving us a fresh new genesis with latests features
eved init $MONIKER1 --chain-id $CHAINID --home $TESTNETS_SHARED_DIR

# key1 - eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn
echo "decorate bright ozone fork gallery riot bus exhaust worth way bone indoor calm squirrel merry zero scheme cotton until shop any excess stage laundry" | eved keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --home $V1 --recover
# key2 - eve1j4rtuq6zm5mmw9xcjmm7gymlj39tvwnt9h4sm2
echo "yard toss ritual ticket dirt address hood stock shiver add client sketch still brave pen win affair orphan employ choose dream sail slogan poverty" | eved keys add $KEY2 --keyring-backend $KEYRING --algo $KEYALGO --home $V2 --recover

# Do via docker?
eved config keyring-backend $KEYRING
eved config chain-id $CHAINID
eved config output "json"


# mkdir $V1 $V2

# Function updates the config based on a jq argument as a string
update_test_genesis () {
    # update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
    cat $TESTNETS_SHARED_DIR/config/genesis.json | jq "$1" > $TESTNETS_SHARED_DIR/tmp_genesis.json && mv $TESTNETS_SHARED_DIR/tmp_genesis.json $TESTNETS_SHARED_DIR/config/genesis.json
}
# Set gas limit in genesis
update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
update_test_genesis '.app_state["gov"]["voting_params"]["voting_period"]="15s"'
# Change chain options to use EXP as the staking denom for craft
update_test_genesis '.app_state["staking"]["params"]["bond_denom"]="ueve"'
update_test_genesis '.app_state["staking"]["params"]["min_commission_rate"]="0.100000000000000000"'
# update from token -> ueve
update_test_genesis '.app_state["mint"]["params"]["mint_denom"]="ueve"'  
update_test_genesis '.app_state["gov"]["deposit_params"]["min_deposit"]=[{"denom": "ueve","amount": "1000000"}]' # 1 eve right now
update_test_genesis '.app_state["crisis"]["constant_fee"]={"denom": "ueve","amount": "1000"}'


# Key 1, copy genesis -> $V1, add gen account, gentx
mkdir -p $V1/config
cp $TESTNETS_SHARED_DIR/config/genesis.json $V1/config/genesis.json

eved add-genesis-account eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn 100000000ueve --keyring-backend $KEYRING --home $V1
eved add-genesis-account eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn 100000000ueve --keyring-backend $KEYRING --home $TESTNETS_SHARED_DIR # also add to main genesis
eved gentx $KEY 1000000ueve --keyring-backend $KEYRING --chain-id $CHAINID --moniker $MONIKER1 -home $V1 --ip 127.0.0.1 #--node-id aaae33661a8286150ad54a512b04bbb96e72b68a -

# Key 2.
mkdir -p $V2/config
cp $TESTNETS_SHARED_DIR/config/genesis.json $V2/config/genesis.json

eved add-genesis-account eve1j4rtuq6zm5mmw9xcjmm7gymlj39tvwnt9h4sm2 100000000ueve --keyring-backend $KEYRING --home $V2
eved add-genesis-account eve1j4rtuq6zm5mmw9xcjmm7gymlj39tvwnt9h4sm2 100000000ueve --keyring-backend $KEYRING --home $TESTNETS_SHARED_DIR
eved gentx $KEY2 1000000ueve --keyring-backend $KEYRING --chain-id $CHAINID --home $V2 --moniker $MONIKER2 --ip 127.0.0.1 #--node-id bbbe33661a8286150ad54a512b04bbb96e72b68a 

# save gentxs back to the root dir
mkdir -p $TESTNETS_SHARED_DIR/config/gentx
cp $V1/config/gentx/*.json $TESTNETS_SHARED_DIR/config/gentx
cp $V2/config/gentx/*.json $TESTNETS_SHARED_DIR/config/gentx

# read -p "Press enter to continue"

eved collect-gentxs --gentx-dir $TESTNETS_SHARED_DIR/config/gentx --home $TESTNETS_SHARED_DIR
eved validate-genesis --home $TESTNETS_SHARED_DIR

# Opens the RPC endpoint to outside connections

# upda the main dirs config toml which we want to persist across both chains values
sed -i '/laddr = "tcp:\/\/127.0.0.1:26657"/c\laddr = "tcp:\/\/0.0.0.0:26657"' $TESTNETS_SHARED_DIR/config/config.toml
sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["\*"\]/g' $TESTNETS_SHARED_DIR/config/config.toml
# copy to the other 2 dirs
mkdir -p $V1/config/ $V2/config/
cp -r $TESTNETS_SHARED_DIR/config/config.toml $V1/config/config.toml
cp -r $TESTNETS_SHARED_DIR/config/config.toml $V2/config/config.toml


# copy main genesis to each dir folder
cp $TESTNETS_SHARED_DIR/config/genesis.json $V1/config/genesis.json
cp $TESTNETS_SHARED_DIR/config/genesis.json $V2/config/genesis.json


# moniker 1 starts a node normally with default values
echo -e "\nStarting the first node ($V1)"
screen -dmS node1 eved start --home $V1 --minimum-gas-prices=0ueve --moniker $MONIKER1 --address "tcp://0.0.0.0:26658" --api.address "tcp://0.0.0.0:1317" --grpc-web.address "0.0.0.0:9091" --grpc.address "0.0.0.0:9090" --p2p.laddr "tcp://127.0.0.1:26656" --rpc.laddr "tcp://127.0.0.1:26657" --proxy_app "tcp://127.0.0.1:26658" --p2p.persistent_peers "bbbe33661a8286150ad54a512b04bbb96e72b68a@127.0.0.1:26667"

echo -e "\nStarting the second node ($V2)"
screen -dmS node2 eved start --home $V2 --minimum-gas-prices=0ueve --moniker $MONIKER2 --address "tcp://0.0.0.0:26668" --api.address "tcp://0.0.0.0:1327" --grpc-web.address "0.0.0.0:9101" --grpc.address "0.0.0.0:9100" --p2p.laddr "tcp://127.0.0.1:26666" --rpc.laddr "tcp://127.0.0.1:26667" --proxy_app "tcp://127.0.0.1:26668" --p2p.persistent_peers "aaae33661a8286150ad54a512b04bbb96e72b68a@127.0.0.1:26657"

# we start in the docker containers here / via compose
# eved start --pruning=nothing  --minimum-gas-prices=0ueve --moniker $MONIKER1 --home $V1