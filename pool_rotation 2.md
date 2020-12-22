# Pool Rotation
[<img src="https://www.bloxstaking.com/wp-content/uploads/2020/04/Blox-Staking_logo_blue.png" width="100">](https://www.bloxstaking.com/)


Public and open staking pool protocols have some challenges around the security of the funds as the protocol is a layer 2 (L2) protocol.
Most notably the creation of new pools and the management of rewards/ penalties of active participants require a dynamic set of participants that can be assigned randomly to pools.

A pool represents an active eth 2.0 validator, identified by a BLS12-381 public key. 
The validator's public and private keys can not be changed once set (at least for now), making any pool rotation require a handoff scheme between pool members at time i and i+1.

The basic key generation utalizes, at it's core, the Joint-Feldman protocol.
A handoff scheme is described as (m,n) -> (m',n')
[An example](https://www.cs.columbia.edu/~wing/publications/Wong-Wing02b.pdf) handoff protocol can bee described as:
1. A set of participants (n) exists with a threshold (m) from a previous DKG round
2. For each ![formula](https://render.githubusercontent.com/render/math?math={i}\in{n}), generate a random polynomial ![formula](https://render.githubusercontent.com/render/math?math=F(x)=c_{i_0})+![formula](https://render.githubusercontent.com/render/math?math=c'_{i_1}X)+..+![formula](https://render.githubusercontent.com/render/math?math=c_{i_(m'-1)}X^(m'-1)) (mod q), where ![formula](https://render.githubusercontent.com/render/math?math=c_{i_0}) is participant's i secret share from time i.
3. Calculate new shares to n' and communicate via secure channels
4. Each ![formula](https://render.githubusercontent.com/render/math?math={i}\in{n'}) can calculate his groups secret share (for time i+1) such: ![formula](https://render.githubusercontent.com/render/math?math=\sum_{1}^{n'}b_i*s'_{i_j}) where ![formula](https://render.githubusercontent.com/render/math?math=b_i) is the lagrange coefficient.
5. After verifying his new share, each ![formula](https://render.githubusercontent.com/render/math?math={i}\in{n'}) will delete his previous outdated share.

## Dynamic Proactive Secret Sharing (DPSS) - Security model
DPSS is a scheme involving a set of participants, together sharing a secret with an (m,n) threshold and dynamic redistribuition of their shares to a (m',n').
As [Ostrovsky and Yung](http://web.cs.ucla.edu/~rafail/PUBLIC/189.pdf) describe in their paper, we assume an attacker might be mobile, meaning he can eventually compromise all participants in the system. As long as t < m of the participants are compromised (either by an attacker or they themselves are malicious) each rotation round, the system is secure.

In a real-word usecase of eth 2.0 staking pools, m will be set to 2/3 as the base protocol defines.
Pool assgnments are random as described below and are built in such a way as to guarantee a very low probability that a malicious party could take over a pool.

#### DPSS - Do shares from differenet epochs can reconstruct the seecret?
If participants get rotated between pools they migh end up having different shares for the same pool from different epochs. Can't they re-construct the secret that way? 
The answer is no for 2 reasons:
1) they only get assigned the shares for their index every time, not different indexes (an index is just the value at the index of the re-distribuition polynomial)
2) shares from different epochs can't re-construct the secret. See [shuffle integrity test](https://github.com/bloxapp/eth2-staking-pools-research/blob/master/go_minimal_pool/crypto/shuffle_integrity_test.go)

## Pool assgniments - min pool size
Each participant is assgined to a pool every epoch randomly to prevent malicious participants to coordinate ahead of time or influence the random source to force their majority in a pool.
With the above in mind, the setup should be very similar to the way eth 2.0 selects committees.
**A good read can be found [here](https://notes.ethereum.org/@vbuterin/rkhCgQteN?type=view#Why-32-ETH-validator-sizes)**

At the heart of this mechanism there is the beacon-chain's random beacon, and some algorithm for deriving the participant's next pool assgniment. eth 2.0 uses the [swap or not algorithm](https://link.springer.com/content/pdf/10.1007%2F978-3-642-32009-5_1.pdf) ([implemention](https://link.springer.com/content/pdf/10.1007%2F978-3-642-32009-5_1.pdf)).
Assuming less than a 1/3 of all participants are malicous, the assginment of participants to pools should guarantee that no single pool has more than 1/3 malicious participants (assuming that < 1/3 are malicous in the selection set).

The way to achieve that is simply by setting the pool size to be big enough, as explained [here](https://notes.ethereum.org/@vbuterin/rkhCgQteN?type=view#Why-32-ETH-validator-sizes).
By binomial distribuition, a minimum pool size of 128 (![formula](https://render.githubusercontent.com/render/math?math=2^7)) will limit the probability of a malicious minority taking over a pool's majority to ![formula](https://render.githubusercontent.com/render/math?math=5.55*10^-15). The actual min pool size is 111 but it was rounded up to the closest exponent of 2.

## [TBD] Pool assgniments - Genesis ceremony
When starting an open protocol for staking pools we should also consider what is the minimal number of participants needed to prevent the small-set attack, a malicous participant could take advantage of low amount of participants (that form a pool) to take over pools failry cheaply. An example would be if a pool has 60 participants, 19 of which already joined. An attacker could complete the necessary 39 and kidnap the group as the selection pool is too small.

The prevent the above, a minimal number of joined participants is needed to kick start the network. exact number **TBD**

## Proactive secret sharing libraries
* CHURP
	* [github](https://github.com/CHURPTeam/CHURP)
	* [paper](https://eprint.iacr.org/2019/017.pdf) 
* [Binance TSS](https://github.com/binance-chain/tss-lib)
* [Unbound blockchain crypto MPC](https://github.com/unbound-tech/blockchain-crypto-mpc)