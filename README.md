# eve

Eve is a modern Cosmos-SDK blockchain intended to push forward the state of the art in Cosmos.  We feel that it will be of particular interest to contract authors who wish to try things that are dangerous or require an adversarial environment for testing.  Eve will be live-fire, so if you think your contract is too hot for Juno or Osmosis, or Neutron, or you wish to try out new features/higher speed, eve will do you well.

## Relationships with other chains and work

* Osmosis style epochs
* Iqlusion style liquid staking
* TGrade style global fees
* Osmosis style TokenFactory with bindings
* Kusama-style level of risk
* e-money style variable block timing

Other intended features:

* Extensive TESTING SUITE

* Airdrop utilities

The fact is, we want lots of this stuff for Juno.  Other chains want other pieces of it.  There's no place where they're packaged together with a stability guarantee, though.
Our success condition here is that anyone can fork eve, and get stable code that's compatible with the rest of Cosmos.

## What is done

* Kusama style risk - DONE
* TGrade global fees - DONE
* Osmosis style TokenFactory - DONE
* Iqlusion style liquid staking - DONE

## NOT DONE YET

* e-money style variable block
* TESTING SUITE (let's be careful plz!)

### TypeScript SDK

Use the [eve TS SDK](https://www.npmjs.com/package/eve-network) for client-side development

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

## Credits

<https://github.com/CosmWasm/token-factory>
<https://github.com/iqlusioninc/liquidity-staking-module>
<https://github.com/cosmos/gaia/tree/main/x/globalfee>

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

<https://twitter.com/gadikian/status/1562019880257277952>

* TBD

## Starting localnet

```bash
make install
make localnet-build
make localnet-start
```
