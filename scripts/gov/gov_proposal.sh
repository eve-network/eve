KEY="test"
CHAINID="eve-testnet"
KEYRING="test"
MONIKER="local-testnet"
KEYALGO="secp256k1"
LOGLEVEL="info"
VALIDATOR="validator"
TESTING_ACCOUNT="vesting_account"
MY_VALIDATOR_ADDRESS=$(eved keys show $VALIDATOR -a --keyring-backend $KEYRING)

eved tx gov submit-legacy-proposal initial-airdrop \
    /home/dangvhb/dangvhbProject/blockchain/notional/eve/scripts/gov/airdrop-proposal.json \
    --title="Test Proposal" \
    --description="testing" \
    --deposit="200000ueve" \
    --fees 200000ueve\
    --from $MY_VALIDATOR_ADDRESS \
    --keyring-backend $KEYRING \
    --chain-id $CHAINID \
    --yes

echo "===> Waiting for transaction to be effective"
sleep 15

echo "===> Voting governance proposal"
eved tx gov vote 1 yes --from $MY_VALIDATOR_ADDRESS --fees 100000ueve --keyring-backend $KEYRING --chain-id $CHAINID --yes

echo "===> Waiting for transaction to be effective"
sleep 15
eved tx bank send $MY_VALIDATOR_ADDRESS eve1258nuq58cz9tfcge5e5egeq69fdvdy7rxmjksa 1000000ueve --fees 100000ueve --keyring-backend $KEYRING --chain-id $CHAINID --yes