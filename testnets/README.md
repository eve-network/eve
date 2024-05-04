# Eve Network Testnet

This testnet will start with the node version `0.0.2`.

## Minimum hardware requirements

- 8-16GB RAM
- 100GB of disk space
- 2 cores

## Genesis Instruction

### Install node

```bash
git clone https://github.com/eve-network/eve.git
cd eve
git checkout v0.0.2
make install
```

### Check Node version

```bash
# Get node version (should be v0.0.2)
eved version

# Get node long version (should be 25e1602d5a29e4ab49addda7e4178b50894999df)
eved version --long | grep commit
```

### Initialize Chain

```bash
rm -rf ~/.eved
eved init develop --chain-id=evenetwork-1
```

### Replace pre-genesis

```bash
# Download the file
curl -s https://raw.githubusercontent.com/eve-network/eve/main/testnets/genesis.json > ~/.eved/config/genesis.json

# Calculate the SHA256 checksum
calculated_checksum=$(shasum -a 256 ~/.eved/config/genesis.json | awk '{ print $1 }')

# Compare with the expected checksum
expected_checksum="244d5a3999dd0851eb338b032a57fbea24a89b4016a7907a9d20c2045c689857"
if [ "$calculated_checksum" = "$expected_checksum" ]; then
    echo "---> Checksum is CORRECT."
else
    echo "---> Checksum is INCORRECT."
fi
```

## Run node

### Setup seeds

```bash
export PERSISTENT_SEEDS="5dd0e206e75e05a21188daffc11440969358fa81@94.130.64.229:26656"
```

### Run node with persistent peers

```bash
eved start --p2p.persistent_peers=$PERSISTENT_SEEDS
```
