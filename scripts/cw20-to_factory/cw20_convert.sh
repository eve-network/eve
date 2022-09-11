
# KEY_NAME="token_info" # this would be in the contract's storage
KEY_NAME="balances"
junod q wasm contract-state raw juno15u3dt79t6sxxa3x3kpkhzsy56edaa5a66wvt3kxmukqjz2sx0hes5sn38g --b64 `echo -n "token_info" | base64` --output json