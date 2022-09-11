# TODO: prefix with owner name so no denoms can conflict with namespace? eveaddr/denom
# TODO: dont allow minting token named eve / ueve?

# = Init token = ($aaa, description, 6 points of percision, 100,000 max supply, maxSupplyNotChangeable, mintable, burnable)
eved tx tokenfactory create-denom "MyToken AAA" aaa 6 100000 false --from eve1 --chain-id eve-t1 --yes -b sync --keyring-backend test
eved q tokenfactory show-denom factory/eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn/aaa
# eved q tokenfactory list-denoms

# = Description =
eved tx tokenfactory update-denom-desc factory/eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn/aaa "new-desc-2 is here for you" --from eve1 --chain-id eve-t1 --yes -b sync --keyring-backend test
# eved q tokenfactory show-denom aaa

# = Token Image =
eved tx tokenfactory update-denom-image factory/eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn/aaa "https://raw.githubusercontent.com/cosmos/chain-registry/master/cosmoshub/images/atom.svg" --from eve1 --chain-id eve-t1 --yes -b sync --keyring-backend test
# eved q tokenfactory show-denom aaa

# = Website =
eved tx tokenfactory update-denom-website factory/eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn/aaa "https://reece.sh" --from eve1 --chain-id eve-t1 --yes -b sync --keyring-backend test
# eved q tokenfactory show-denom aaa



# = Change Max Supply = # is this even needed?
# 'Cannot revert change maxsupply' if false in create-denom
eved tx tokenfactory update-denom-change-supply factory/eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn/aaa true --from eve1 --chain-id eve-t1 --yes -b sync --keyring-backend test
# eved q tokenfactory show-denom aaa


# = Mint Tokens -> An Account =
# key found in ./test_node.sh
eved tx tokenfactory mint-and-send-tokens factory/eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn/aaa 9 eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn --from eve1 --chain-id eve-t1 --yes -b sync --keyring-backend test
eved q bank balances eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn --denom=factory/eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn/aaa


# = Transfer Ownership of a denom = (do we want this? not very CW20 like)
eved tx tokenfactory update-owner factory/eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn/aaa eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn --from eve1 --chain-id eve-t1 --yes -b sync --keyring-backend test