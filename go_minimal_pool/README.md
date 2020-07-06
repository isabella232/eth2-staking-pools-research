# ETH 2.0 Go minimal staking pool implementation - Research
[<img src="https://www.bloxstaking.com/wp-content/uploads/2020/04/Blox-Staking_logo_blue.png" width="100">](https://www.bloxstaking.com/)

## This is still under development!!

A basic and minimal POC for staking pools for eth 2.0 based on this [research](https://github.com/bloxapp/eth2-staking-pools-research).

### What it does?
* Initial DKG (without VSS!)
* contructs epochs and rotates participants randomly between them
* during a rotation it does a redistribution of shares very naively at the moment (again, no VSS).
* It has no netwokring, all participants send messages via function calls.

This project is a result of the [python_minimal_pool](https://github.com/bloxapp/eth2-staking-pools-research/tree/master/python_minimal_pool). It was too slow for pairing operations.