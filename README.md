# eve

Eve is a modern Cosmos-SDK blockchain intended to push forward the state of the art in cosmos.  We feel that it will be of particular interest to contract authors who wish to try things that are dangerous or require an adversarial environment for testing.  Eve will be live-fire, so if you think your contract is too hot for Juno or Osmosis, or Neutron, or you wish to try out new features/higher speed, eve will do you well. 

## Relationships to other chains and work

* Osmosis-style epochs
* kusama-style level of risk
* Tgrade style global fees
* e-money style variable block timing
* suport for iqlusion style liquid staking

Other intended features:
* Extensive TESTING SUITE
* Airdrop utilities

Fact is, we want lots of this stuff for Juno.  Other chains want other pieces of it.  There's no place where they're packaged together with a stability guarantee, though.

Our success condition here is that anyone can fork eve, and get stable code that's compatible with the rest of cosmos.

#### What is done:
* Osmosis-style epochs - DONE
* Kusama-style risk - DONE
* Tgrade global fees - DONE
* iqlusion style liquid staking - DONE

#### NOT DONE YET!
* e-money style variable block
* TESTING SUITE (lets be careful plz!)


### Upstreaming

In call cases, when stuff changes here and we find it nice, we'll feed it back to parent repos.

### Downstreaming

Likey stuff built and tested here feeds into Juno and a number of other chains.


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

  
  

* Pebbledb by default
* cosmwasm
* token factory
* variable length blocks made on demand with a 3-minute heartbeat
* genesis validator set limited to contributors

## Use of Eve as a framework

Numerous ideas in eve came about as a result of work on Craft Economy and Juno.  We were also inspired by recent moves to make the Osmosis epoch module importable-- it is better than the one in the SDK, but isn't importable at this time. Thus, Eve herself is becoming an sdk or framework for light, efficient CosmWasm Chains that natively support liquid staking.  

Features developed in eve are designed to be importable directly from Eve's repository, and eve won't have her own SDK fork, though there may be a fork of wasmd. 

To use eve as a framework to make your own chain, just fork eve. 

## Credits
https://github.com/clockworkgr/tokenfactory  (v45 -> v46.1)
https://github.com/iqlusioninc/liquidity-staking-module
WasmBindings - Fork of Osmosis's wasmbindings

## Imperative

From eve will flow new designs and techniques that increase development velocity and ease of use across Cosmos.  


## Contributing

* Contributors will be added to this repository with write access.
* Chronic contributors will be added as owners
* We'll move this out of notional-labs sooner or later and owners will own the eve-network (or such) GitHub org
* you can use https://gitpod.io to contribute easily without setting up a development environment
* you contribute here, in this repository, and there are many ways to contribute
  * code == docs
  * docs == code
  * logo == code
  * docs == logo


## Airdrop


https://twitter.com/gadikian/status/1562019880257277952

* TBD

## Starting localnet

```
make install
make localnet-build
make localnet-start
```