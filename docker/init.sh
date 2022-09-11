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

V1=$SCRIPT_DIR/.testnets/v1 # /home/reece/Desktop/Programming/Go/eve/docker/.testnets/v1
V2=$SCRIPT_DIR/.testnets/v2 # /home/reece/Desktop/Programming/Go/eve/docker/.testnets/v2

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

# Function updates the config based on a jq argument as a string
update_test_genesis () {
    # update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
    cat $TESTNETS_SHARED_DIR/config/genesis.json | jq "$1" > $TESTNETS_SHARED_DIR/tmp_genesis.json && mv $TESTNETS_SHARED_DIR/tmp_genesis.json $TESTNETS_SHARED_DIR/config/genesis.json
}

# get current time in bash which matches teh ISO8601 format 2022-09-11T22:22:58.150405469Z
seconds_in_the_future=`expr $(date +%S)` # 10 seconds in the future
if [ $(date +%S) -gt 30 ]; then    
    if [ $(date +%S) -lt 45 ]; then
        seconds_in_the_future=`expr $(date +%S) + 5`
    fi
else 
    seconds_in_the_future=`expr $seconds_in_the_future + 5`
fi

# the_time=$(date +%FT%H:%M:${seconds_in_the_future}.000000000Z --utc)
the_time=$(date +%FT%H:%M:${seconds_in_the_future}Z --utc)
echo "Setting genesis time to $the_time"
# exit 0
# exit 0

# Set gas limit in genesis
update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
update_test_genesis `printf '.genesis_time="%s"' $the_time`
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

eved gentx $KEY 1000000ueve --keyring-backend $KEYRING --chain-id $CHAINID --home $V1 --moniker $MONIKER1 --commission-rate=0.10 --commission-max-rate=1.0 --commission-max-change-rate=0.01 --min-self-delegation "1" --ip 127.0.0.1 --node-id 017e3b0c9a050091cc2bc609af9fb861d3710215

# Key 2.
mkdir -p $V2/config
cp $TESTNETS_SHARED_DIR/config/genesis.json $V2/config/genesis.json

eved add-genesis-account eve1j4rtuq6zm5mmw9xcjmm7gymlj39tvwnt9h4sm2 100000000ueve --keyring-backend $KEYRING --home $V2
eved add-genesis-account eve1j4rtuq6zm5mmw9xcjmm7gymlj39tvwnt9h4sm2 100000000ueve --keyring-backend $KEYRING --home $TESTNETS_SHARED_DIR

eved gentx $KEY2 1000000ueve --keyring-backend $KEYRING --chain-id $CHAINID --home $V2 --moniker $MONIKER2 --commission-rate=0.10 --commission-max-rate=1.0 --commission-max-change-rate=0.01 --min-self-delegation "1" --ip 127.0.0.1 --node-id 017e3b0c9a050091cc2bc609af9fb861d3710215

# KEY1_NODE_ID=`eved tendermint show-node-id --home $V1`
# KEY2_NODE_ID=`eved tendermint show-node-id --home $V2`

# save gentxs back to the root dir
mkdir -p $TESTNETS_SHARED_DIR/config/gentx
mv $V1/config/gentx/*.json $V1/config/gentx/v1_node.json

cp $V1/config/gentx/*.json $TESTNETS_SHARED_DIR/config/gentx
cp $V2/config/gentx/*.json $TESTNETS_SHARED_DIR/config/gentx

# read -p "Press enter to continue"

eved collect-gentxs --gentx-dir $TESTNETS_SHARED_DIR/config/gentx --home $TESTNETS_SHARED_DIR
eved validate-genesis --home $TESTNETS_SHARED_DIR

# Opens the RPC endpoint to outside connections

# upda the main dirs config toml which we want to persist across both chains values
sed -i '/laddr = "tcp:\/\/127.0.0.1:26657"/c\laddr = "tcp:\/\/0.0.0.0:26657"' $TESTNETS_SHARED_DIR/config/config.toml
sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["\*"\]/g' $TESTNETS_SHARED_DIR/config/config.toml

# edit just V1 so it has a unique rosetta address
sed -i 's/address = ":8080"/address = ":8079"/g' $V2/config/app.toml

# copy to the other 2 dirs
mkdir -p $V1/config/ $V2/config/
cp -r $TESTNETS_SHARED_DIR/config/config.toml $V1/config/config.toml
cp -r $TESTNETS_SHARED_DIR/config/config.toml $V2/config/config.toml


# copy main genesis to each dir folder
cp $TESTNETS_SHARED_DIR/config/genesis.json $V1/config/genesis.json
cp $TESTNETS_SHARED_DIR/config/genesis.json $V2/config/genesis.json

read -p "Press enter to continue"

# goal is to do this via each docker container. I guess we could support screens too if someone wanted that
# 6:04PM INF Error reconnecting to peer. Trying again addr={"id":"aaae33661a8286150ad54a512b04bbb96e72b68a","ip":"127.0.0.1","port":26657} err="auth failure: secret conn failed: proto: BytesValue: wiretype end group for non-group" module=p2p tries=0
echo -e "\nStarting the first node ($V1)"
screen -dmS n1 eved start --home $V1 --moniker $MONIKER1 --address "tcp://0.0.0.0:26658" --api.address "tcp://0.0.0.0:1317" --grpc-web.address "0.0.0.0:9091" --grpc.address "0.0.0.0:9090" --p2p.laddr "tcp://127.0.0.1:26656" --rpc.laddr "tcp://127.0.0.1:26657" --proxy_app "tcp://127.0.0.1:26658" --p2p.persistent_peers "017e3b0c9a050091cc2bc609af9fb861d3710215@127.0.0.1:26666"

echo -e "\nStarting the second node ($V2)"
screen -dmS n2 eved start --home $V2 --moniker $MONIKER2 --address "tcp://0.0.0.0:26668" --api.address "tcp://0.0.0.0:1327" --grpc-web.address "0.0.0.0:9101" --grpc.address "0.0.0.0:9100" --p2p.laddr "tcp://127.0.0.1:26666" --rpc.laddr "tcp://127.0.0.1:26667" --proxy_app "tcp://127.0.0.1:26668" --p2p.persistent_peers "017e3b0c9a050091cc2bc609af9fb861d3710215@127.0.0.1:26656"