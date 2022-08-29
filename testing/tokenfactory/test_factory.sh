# http://localhost:1317/clockworkgr/tokenfactory/tokenfactory/denom

# TODO: denom & ticker = the same?
# TODO:prefix with owner name so no denoms can conflict with namespace? eveaddr/denom

eved tx tokenfactory create-denom reece "my denom" 6 100000 false --from eve1 --chain-id eve-t1 --yes -b sync --keyring-backend test

eved q tokenfactory list-denom

eved q tokenfactory show-denom reece


# update denom - max-supply & can change should be optional
# eved tx tokenfactory update-denom reece "new-desc-2" 100001 false --from eve1 --chain-id eve-t1 --yes -b sync --keyring-backend test
eved tx tokenfactory update-denom-desc reece "new-desc-2" --from eve1 --chain-id eve-t1 --yes -b sync --keyring-backend test

# add:
# update-denom-url ?
# update-denom-supply (increase/descrease, and toggle change to false) IF they are allowed too.

# todo: dont allow minting token named eve / ueve?

# mint tokens to eve1 account & check that they have the tokens now.
eved tx tokenfactory mint-and-send-tokens reece 9 eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn --from eve1 --chain-id eve-t1 --yes -b sync --keyring-backend test
eved q bank balances eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn --denom=reece