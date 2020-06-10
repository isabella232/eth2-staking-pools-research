# ETH 2.0 Decentralized Staking Pools - Research

This repo aims to have in one place all the research around decentralized staking pools.

- Distribuited key generation and redistribuition
- Rewards/ penalties
- Consensus
- bilinear pairings and BLS12-381 keys
- Networking  

### Overview
The backbone of a decentralized staking pool is in distribuiting the control of the keys that control the validator and its withdrawal key. You can think of it as a giant multisig setup with some N-of-M threshold for signing attestations, block proposals and withddrawal transactions.
A good starting point could be [this](https://www.youtube.com/watch?v=Jtz9b7yWbLo) presentation.

If we add a consensus protocol that rewards and punishes pool participants, controls withdrawal and onboarding processes then we have a full mechanism for decentralzied staking pools network.
One issue that raises immediatley is the security around the oboarding process, how can we guarantee that a formed pool will includee (worst case) no more that 1/3 malicious participants?
According to Binomial distribuition we can calculate how big does a pooll needs to be so, this is similar to ethereum committee selection as explained [here](https://notes.ethereum.org/@vbuterin/rkhCgQteN?type=view#Why-32-ETH-validator-sizes). 
Another issue is how big the actual set of participants to select from, if it's too small a malicious participant can hijack entire pools, **example**: a pool consists of 128 participants but the available set of participants is only of 40 participants. A malicious opponent can quickly add 88 of his own participants and jsut steal the other participants eth. 
This scenario is very possible as there is no guarantee that a large pool of non allocated participants will exist, and even if it is then it will provide bad experience for users to just wait.

A solution to this porblem is a continous rotation of all participants between the pools such that new created pools can share the set of "rotating pariticipants" to create the necessary randomness in allocating new validators.
The requirements from individual participants are as they are for a full 32 ETH valdiator:
- Be online and execute your duties
- Do not execute on a duty that will cause slashing

A participant that does not execute on his assigned pool duties or signs on a duty that can cause slashing will get penalized.
A participant that executes on his duties will get rerwarded.

