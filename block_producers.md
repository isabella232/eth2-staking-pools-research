# Consensus - Block Producers
[<img src="https://www.bloxstaking.com/wp-content/uploads/2020/04/Blox-Staking_logo_blue.png" width="100">](https://www.bloxstaking.com/)

Block producers(BP) have the responsibility of commiting a block to a chain containing the following data:
- epoch's state root
- [new pool creations](https://github.com/bloxapp/eth2-staking-pools-research/blob/master/new_pools.md)
- [participant/ BP slashing](https://github.com/bloxapp/eth2-staking-pools-research/blob/master/pool_duties.md)
- eth withdrawal requests
- CDT withdrawal requests

For their work, block producers get rewards in CDT.

## Block Producer - Rewards/ Penalties/ Slashing
If a BP executes his duties corretcly he will get rewarded otherwise he will be peenalaized or slashed.

BP receives a reward for corrctly commiting a new block to the chain.

BP will get penalized for not submitting a valid block on his turn
BP will get penalized for submitting a valid block that did not get enough votes (not finalized)

BP will get slashed for submitting similar blocks (same block height with different data)
BP will get slashed if his balance get's below a certain amount (TBD)

## Becoming a Block Producer
Becoming a block producer requires staking CDT in an eth1 contract. 
The staked CDT acts as collateral for the BP's actions. 

## Recieving Rewards
As the BP gains CDT rewards on the pool chain (seperate from eth) he can decide to withdraw them back to eth (for the reasoninng behind keeping CDT on eth see cdt2 page). `
withdrawal is a request that can be submitted by any BP on the pools chain (as a block data argument), once that request (and the block it's in) is finalized the BP can submit a request to eth1 with the specific pool chain block hash and the amount to withdraw.

There is a waiting period in which any other BP can dispute this request (on eth1) and potetially slash the requesting BP and get rewarded. 
In plain words, a BP can submit a withdrawal request from the pool chain to eth1 and if nobody disputes it he will get his CDT on eth1. If a valid dispute is submitted he will lose his dposit and any remaining CDT on thee pool chain.

## Fee Rewards
As the BPs are the backbone of pool chain they are entitled for recieving any fees created by the network. 
Fees might include: user deposit/ withdrawal, slashed users and so on.