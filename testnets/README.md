# Eve Testnet

This testnet will start with the node version `v0.0.3`.

## Minimum hardware requirements

- 8-16GB RAM
- 100GB of disk space
- 2 cores

## Genesis Instruction

### Install node

```bash
git clone https://github.com/eve-network/eve.git
cd eve
git checkout v0.0.3
make install
```

### Check Node version

```bash
# Get node version (should be v0.0.3)
eved version

# Get node long version (should be 1f0f1f82a8225b23341bbabd2a034ce7415d7e3d)
eved version --long | grep commit
```

### Initialize Chain

```bash
eved init MONIKER --chain-id=evenetwork-1
```

### Download pre-genesis

```bash
curl -s https://raw.githubusercontent.com/eve-network/eve/main/testnets/pre_genesis.json > ~/.eved/config/genesis.json
```

## Create gentx

Create wallet

```bash
eved keys add KEY_NAME
```

Fund yourself `1000000000ueve`

```bash
eved genesis add-genesis-account $(eved keys show KEY_NAME -a) 1000000000ueve
```

Use half (`1000000ueve`) for self-delegation

```bash
eved genesis gentx KEY_NAME 1000000ueve --chain-id=evenetwork-1
```

If all goes well, you will see a message similar to the following:

```bash
Genesis transaction written to "/home/user/.eved/config/gentx/gentx-******.json"
```

### Submit genesis transaction

- Fork this repo
- Copy the generated gentx json file to `testnets/gentx/`
- Commit and push to your repo
- Create a PR on this repo
