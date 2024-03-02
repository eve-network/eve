# eve

Eve is a governance minimized CosmWasm blockchain intended to push forward the state of the art in Cosmos.  

Contract uploads are not permissioned in any way, and governance thresholds have been adjusted to ensure that contract authors are not disturbed in their work.   This chain emphasizes decentralization, interchain capability and governance minimization.  

Eve won't compete with your contracts, and views herself as a beautiful, fertile platform for authors to bring their visions to life. 


## Software at launch

Note: it is possible that we fall back to v0.47.x if simply avoiding vote extensions proves to be ineffective in stopping the issues seen in deployments of comet v0.38.x

* cosmos-sdk v0.50.x without vote extensions, with cyber-sdk gpu support
* cometbft v0.38.x without vote extensions
* CosmWasm v0.50.x
* WasmVM v1.5.x
* IBC-go v8.1.x


## Planned software

* interchain security provider
* mesh security provider

### Upstreaming

In all cases, when stuff changes here and we find it nice, we'll feed it back to parent repositories.

### Downstream

Likely stuff built and tested here feeds into Juno and several other chains.

## Team

The individuals and groups in the git history are eve's team. 

## No Promises

Eve doesn't make promises and may not be exactly as described because she's an exploration of ideas first and foremost.  This includes the airdrop.  

## Economic Information

Anything described below can be modified by making a pull request.  

* genesis:  ????????
  * embedded in binary after launch
  * no "dev fund"
  * community tax set to 20% initially
  * 1 validator = 1 eve
  * 10% minimum commission
  * 50 validators

* PebbleDB by default
* Cosmwasm
* Token Factory
* Variable length blocks made on demand with a 3-minute heartbeat
* Genesis validator set limited to contributors

## Use of Eve as a framework

Numerous ideas in eve came about as a result of work on Craft Economy and Juno.  We were also inspired by recent moves to make the Osmosis epoch module importable-- it is better than the one in the SDK but isn't importable at this time. Thus, Eve herself is becoming an SDK or framework for light, efficient CosmWasm Chains that natively support liquid staking.  

Features developed in eve are designed to be importable directly from Eve's repository, and eve won't have her SDK fork, though there may be a fork of Wasmd.

To use eve as a framework to make your chain, just fork eve.


## Imperative

From eve will flow new designs and techniques that increase development velocity and ease of use across Cosmos.  

## Contributing

* Contributors will be added to this repository with write access.
* Chronic contributors will be added as owners
* We'll move this out of notional-labs sooner or later and owners will own the eve-network (or such) GitHub org
* you can use <https://gitpod.io> to contribute easily without setting up a development environment
* you contribute here, in this repository, and there are many ways to contribute
  * code == docs
  * docs == code
  * logo == code
  * docs == logo

## Airdrop

Nothing is guaranteed here, but whales won't be pruned.  Conversely, non-productive validators and their delegators may be pruned. 

