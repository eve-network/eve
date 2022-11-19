KEY="test"
CHAINID="eve-testnet"
KEYRING="test"
MONIKER="local-testnet"
KEYALGO="secp256k1"
LOGLEVEL="info"
VALIDATOR="validator"
TESTING_ACCOUNT="vesting_account"
MY_VALIDATOR_ADDRESS=$(eved keys show $VALIDATOR -a --keyring-backend $KEYRING)
echo $MY_VALIDATOR_ADDRESS
eved tx gov submit-legacy-proposal initial-airdrop \
    /home/dangvhb/dangvhbProject/blockchain/notional/eve/scripts/gov/airdrop-proposal.json \
    --title="Test Proposal" \
    --description="testing" \
    --deposit="1000ueve" \
    --fees 1000ueve\
    --from $MY_VALIDATOR_ADDRESS \
    --keyring-backend $KEYRING \
    --chain-id $CHAINID