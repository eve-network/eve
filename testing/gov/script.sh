# Submit a change from default fees (none) to some fee
eved tx gov submit-legacy-proposal param-change testing/gov/fee_param.json --from eve1 --keyring-backend test --chain-id eve-t1 --yes
eved tx gov vote 1 yes --from eve1 --keyring-backend test --chain-id eve-t1 --yes

eved q globalfee minimum-gas-prices --output json

eved tx bank send eve1 eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn 1ueve --keyring-backend test --chain-id eve-t1 --gas 100000 --fees 250ueve --yes
eved tx bank send eve1 eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn 1ueve --keyring-backend test --chain-id eve-t1 --gas 200000 --fees 499ueve --yes # fail, should be 500


# set fee to 0, run all these of these at the same time.
SEQ=$(eved q account eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn --output json | jq -r '.sequence')
LAST_PROP_ID=$(eved q gov proposals --output json | jq -r '.proposals | last | .id')
eved tx gov submit-legacy-proposal param-change testing/gov/lower_min_fee.json --from eve1 --keyring-backend test --chain-id eve-t1 --yes --fees 500ueve --sequence $SEQ
eved tx gov vote $((LAST_PROP_ID+1)) yes --from eve1 --keyring-backend test --chain-id eve-t1 --yes --fees 500ueve --sequence $((SEQ+1))


eved q globalfee minimum-gas-prices --output json

eved tx bank send eve1 eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn 1ueve --keyring-backend test --chain-id eve-t1 --gas 100000 --fees 0ueve --yes
eved tx bank send eve1 eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn 1ueve --keyring-backend test --chain-id eve-t1 --gas 200000 --fees 1ueve --yes