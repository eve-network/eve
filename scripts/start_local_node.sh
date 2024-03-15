#!/bin/bash
set -eu

EVE_HOME=~/.eve-local
EVED="eved --home ${EVE_HOME}"
CHAIN_ID=eve-local-1
DENOM=ueve

EVE_ADMIN_MNEMONIC="gorilla bind ghost erode play crack ancient flight mountain floor rent flip lens split today winter oil arctic recycle lab reform habit keep cactus"
EVE_VAL_MNEMONIC="jeans agree enter oak sure amateur ride ceiling museum bunker weekend fruit give truth blush lucky ball chunk regret mirror leader pudding mirror web"

MAX_DEPOSIT_PERIOD="20s"
VOTING_PERIOD="20s"
UNBONDING_TIME="20s"

config_toml="${EVE_HOME}/config/config.toml"
client_toml="${EVE_HOME}/config/client.toml"
app_toml="${EVE_HOME}/config/app.toml"
genesis_json="${EVE_HOME}/config/genesis.json"

rm -rf ${EVE_HOME}

$EVED init eve-local --chain-id $CHAIN_ID

sed -i -E "s|minimum-gas-prices = \".*\"|minimum-gas-prices = \"0${DENOM}\"|g" $app_toml
sed -i -E '/\[api\]/,/^enable = .*$/ s/^enable = .*$/enable = true/' $app_toml

sed -i -E "s|chain-id = \"\"|chain-id = \"${CHAIN_ID}\"|g" $client_toml
sed -i -E "s|keyring-backend = \"os\"|keyring-backend = \"test\"|g" $client_toml
sed -i -E "s|node = \".*\"|node = \"tcp://localhost:26657\"|g" $client_toml

sed -i -E "s|\"stake\"|\"${DENOM}\"|g" $genesis_json

jq '.app_state.gov.params.max_deposit_period = $newVal' --arg newVal "$MAX_DEPOSIT_PERIOD" $genesis_json > json.tmp && mv json.tmp $genesis_json
jq '.app_state.gov.params.voting_period = $newVal' --arg newVal "$VOTING_PERIOD" $genesis_json > json.tmp && mv json.tmp $genesis_json

# hack since add-comsumer-section is built for dockernet
rm -rf ~/.eve-loca1
cp -r ${EVE_HOME} ~/.eve-loca1

$EVED add-consumer-section 1
jq '.app_state.ccvconsumer.params.unbonding_period = $newVal' --arg newVal "$UNBONDING_TIME" $genesis_json > json.tmp && mv json.tmp $genesis_json

rm -rf ~/.eve-loca1

# Start the daemon in the background
$EVED start
#pid=$!
#sleep 10
#
## Add a governator
#echo "Adding governator..."
#pub_key=$($EVED tendermint show-validator)
#$EVED tx staking create-validator --amount 1000000000${DENOM} --from val \
#    --pubkey=$pub_key --commission-rate="0.10" --commission-max-rate="0.20" \
#    --commission-max-change-rate="0.01" --min-self-delegation="1" -y
#
## Bring the daemon back to the foreground
#wait $pid
