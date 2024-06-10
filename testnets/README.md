# Eve Testnet

This testnet will start with the node version `v0.0.4`.

## Minimum hardware requirements

- 8-16GB RAM
- 100GB of disk space
- 2 cores

## Genesis Instruction

### Install node

```bash
git clone https://github.com/eve-network/eve.git
cd eve
git checkout v0.0.4
make install
```

### Check Node version

```bash
# Get node version (should be v0.0.4)
eved version

# Get node long version (should be f31035e45725505a001e4f45a7fdd85be22e110a)
eved version --long | grep commit
```

### Initialize Chain

```bash
eved init MONIKER --chain-id=evenetwork-1
```

### Download genesis

```bash
curl -s https://raw.githubusercontent.com/eve-network/eve/main/testnets/genesis.json > ~/.eved/config/genesis.json
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


## Setup validator
Set up your node and add seed-nodes

    bc7506bbcf9fef7182f3a4628f9fc62c25c2fd27@94.130.64.229:26656
