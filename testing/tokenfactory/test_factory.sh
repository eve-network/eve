SUBDENOM="reece1"
COIN_NAME="factory/eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn/$SUBDENOM"
TX_PARAMS="--from eve1 --chain-id eve-t1 --yes --broadcast-mode sync --keyring-backend test"


SEQ=$(eved q account eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn --output json | jq -r '.sequence')

# Requires 1eve to make.
eved tx tokenfactory create-denom $SUBDENOM $TX_PARAMS --sequence $SEQ

eved q tokenfactory denom-authority-metadata $COIN_NAME

eved tx tokenfactory mint 100$COIN_NAME $TX_PARAMS --sequence $((SEQ+1))


eved q bank balances eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn --denom $COIN_NAME

eved tx tokenfactory burn 5$COIN_NAME $TX_PARAMS --sequence $((SEQ+2))

# eved tx tokenfactory change-admin $COIN_NAME eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn $TX_PARAMS

eved q tokenfactory denoms-from-creator eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn 