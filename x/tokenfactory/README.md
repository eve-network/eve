### TokenFactory

Design considerations:
- Osmosis
- CW20

// name, Denom, decimals, total_supply, canChangeTotalSupply
// extra: description, URL, website, token Image, Twitter? (metadata message)

```
CW20:
- name, symbol, decimals, total_supply, mint (Optional<MinterData>)

Osmosis:
- factory/{creator address}/{subdenom}

- denoms = [a-zA-Z0-9./] 
other: https://github.com/osmosis-labs/osmosis/tree/main/x/tokenfactory

admin:
- mint denom to any account
- burn denom from any account (do we want an admin to get to do this? We should maybe only allow minting. Then the holder of a token can burn it if they wanted too.)
- transfer denom between any 2 accounts

- denomMetadata
```



Add to spec:
```
TODO: (anyone who has the token can send them in to be burned)
burn [amount] [flags] - only allow burning from your own address? or enable this as a toggleable option on create
```

Denom format: (query)
```
denom:
  canChangeMaxSupply: false
  denom: aaa
  description: new-desc-2 is here for you
  maxSupply: 100000
  name: MyToken AAA
  owner: eve1hj5fveer5cjtn4wd6wstzugjfdxzl0xpysfwwn
  precision: 6
  supply: 9
  token_image: https://raw.githubusercontent.com/cosmos/chain-registry/master/cosmoshub/images/atom.svg
  website: https://reece.sh
```