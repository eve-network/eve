```
eved keys add [keyname] --recover
```

# Validator setup instructions
## Genesis Tx:
```bash
# Validator variables
KEYNAME='reece' # your [keyname]
MONIKER='pbcups_testnet'
SECURITY_CONTACT="reece@reece.sh"
WEBSITE="https://reece.sh"
MAX_RATE='0.20'        # 20%
COMMISSION_RATE='0.05' # 5%
MAX_CHANGE='0.05'      # 5%
CHAIN_ID='eve-tn1'
PROJECT_HOME="${HOME}/.eved"
KEYNAME_ADDR=$(eved keys show $KEYNAME -a)
# /Validator variables

# echo -e "$KEYNAME\n$MONIKER\n$DETAILS\n$SECURITY_CONTACT\n$WEBSITE\n$MAX_RATE\n$COMMISSION_RATE\n$MAX_CHANGE\n$CHAIN_ID\n$HOME_DIR\n$KEYNAME_ADDR"

# Remove old files if they exist
eved tendermint unsafe-reset-all
rm $HOME/.eved/config/genesis.json
rm $HOME/.eved/config/gentx/*.json

# Give yourself 1exp for the genesis Tx signed
eved init "$MONIKER" --chain-id $CHAIN_ID --staking-bond-denom ueve
eved add-genesis-account $KEYNAME_ADDR 1000000ueve

# genesis transaction using all above variables
eved gentx $KEYNAME 1000000ueve \
    --home=$PROJECT_HOME \
    --chain-id=$CHAIN_ID \
    --moniker="$MONIKER" \
     --commission-max-change-rate=$MAX_CHANGE \
    --commission-max-rate=$MAX_RATE \
    --commission-rate=$COMMISSION_RATE \
    --security-contact=$SECURITY_CONTACT \
    --website=$WEBSITE \
    --details=""

# Get that gentx data easily -> your home directory
cat ${PROJECT_HOME}/config/gentx/gentx-*.json

# Download the file from $HOME/MONIKER.json & upload to the discord channel
echo -e "\n\n\nPlease paste the contents above into the discord channel (( ${PROJECT_HOME}/config/gentx/gentx-*.json ))"
echo -e "     (also remember to backup ~/.craftd/node_key.json && ~/.craftd/priv_validator_key.json)\n"
echo -e "     Your peer: `echo $(craftd tendermint show-node-id)@$(curl -s ifconfig.me):26656`"
# DO NOT SUBMIT A PR FOR THIS, POST IN THE DISCORD `gentx-submit` CHANNEL THANK YOU!
```