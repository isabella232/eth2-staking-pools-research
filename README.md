[<img src="./img/header.png" width="1000">](https://www.bloxstaking.com/)

This repo aims to have in one place all the research around decentralized staking pools for eth 2.0.

### ETH 2.0 Decentralized Staking Pools - Summary
The backbone of decentralized staking pools is in distributing  the control of the keys that control the validator and its withdrawal key. You can think of it as a giant multi-signature setup with some M-of-N threshold for signing attestations, block proposals and withdrawal transactions.\
A good starting point could be [this](https://www.youtube.com/watch?v=Jtz9b7yWbLo) presentation.

Adding a consensus protocol that rewards and punishes pool participants, controls withdrawal and on-boarding then we have a full protocol for an open decentralized staking pools network.\
The key word here is open as in autonomous and open to join by anyone.

### Overview
Conceptually a pools network can be thought of as a 3 layer stack. 

<b>Layer 1</b> - Every pool is a collection of 32 ETH validators represented by a public BLS key.\
<b>Layer 2</b> - Every pool's public key has a corresponding distributed  private key controlled via an SSV group of operators.\
<b>Layer 3</b> - All pools (and their SSV groups) coordinate via a consensus layer which also deals with rewards, penalties, slashing, creation/ liquidation of pools and more.\

<div style="text-align:center"><img src="./img/design.png" width="400"></div>
<sub>Pools network high-level architecture.</sub><br /><br />


The network has 2 actors: block producers (BP) and staker.\
A BP is a bonded actor (staked)  which has the responsibility of executing top consensus (attest to blocks, propose blocks and more) and local SSV group assignments (mostly eth2 network duties).<br /><br />
A block producer is economically incentivized to run a pool node, participate in the network and more. The block producer's collateral is also staked.\
<i>For more information regarding network economics click [here](https://www.bloxstaking.com/blog/ethereum-2-0/an-introduction-to-decentralized-staking-pools/)</i>.<br /><br />

A staker that deposited ETH into a smart contact to stake in a pool will, in return, mint a fungible ERC-20 token representing his stake + future rewards. At deposit time the amount minted of that token will be 1:1 to the deposit amount, as time goes on and the network as a whole gains reward the user's ERC-20 token balance (and so for every other user) will grow relatively.\
The ERC-20 token creates instant transferability, detached from a specific pool liquidation.<br /><br />

Pool liquidation is an event that is triggered by an ERC-20 token holders specifically requesting to convert back to ETH. This is also dependent on phase 1.5 of eth2 to complete, see [discord discussion](https://discord.com/channels/595666850260713488/748530848470663208/748533989433802772)



### Research
* [Pools mini paper + CDT2.0](pools_mini_paper.pdf) for in-depth details
* [DKG theory](https://github.com/bloxapp/eth2-staking-pools-research/blob/master/dkg.md) 
* [DKG + key rotation POC](https://github.com/bloxapp/eth2-staking-pools-research/tree/master/go_minimal_pool)
* SSV nodes
    * [python ibft implementation](https://github.com/dankrad/python-ibft)
    * [python SSV node](https://github.com/dankrad/python-ssv)
    * [SSV compatible valdiator](https://github.com/alonmuroch/prysm/tree/ssv)
* [BLS keys](https://medium.com/@alonmuroch_65570/bls-signatures-part-1-overview-47d9eebf1c75) + [bi-linear pairings](https://medium.com/@alonmuroch_65570/bls-signatures-part-2-key-concepts-of-pairings-27a8a9533d0c)
* [Isolated Casper + Ghost consensus as candidate for the pool's network consensus](https://github.com/bloxapp/go-casper-ghost-SDK)
* [Network tokonomics](https://www.bloxstaking.com/blog/ethereum-2-0/an-introduction-to-decentralized-staking-pools/)

