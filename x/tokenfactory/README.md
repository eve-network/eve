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
- changeAdmin

- denomMetadata
```



This token factory spec (WIP)
```
Txs:
create-denom [denom] [flags]
mint [amount] [address] [flags]
change-admin [denom] [new-admin-address] [flags]

burn [amount] [flags] - only allow burning from your own address? or enable this as a toggleable option on create


queries:
- metadata [denom] 
- denoms-from-creator [creator address] ?
- others which are already there


- list-denoms
- show-denom [denom]
```