# Pool Rotation
Public and open staking pool protocols have some challenges around the security of the funds as the protocol is an layere 2 (L2) protocol.
Most notably the creation of new pools and the management of rewards/ penalties of active participants require a dynamic set of participants that can be assigned randomly to pools.

A pool represents an active eth 2.0 validator, identified by a BLS12-381 public key. 
The validator's public and private keys can not be changed once set (at least for now), making any pool rotation require a handoff scheme between pool members at time i and i+1.

The basic key generation utalizes, at it's core, the Joint-Feldman protocol.
A handoff scheme is described as (m,n) -> (m',n')
[A naive](https://www.cs.columbia.edu/~wing/publications/Wong-Wing02b.pdf) handoff protocol can bee described as:
1. A set of participants (n) exists with a threshold (m) exists from a previous DKG round
2. For each ![formula](https://render.githubusercontent.com/render/math?math={i}\in{n}), generate a random polynomial ![formula](https://render.githubusercontent.com/render/math?math=F(x)=c_{i_0})+![formula](https://render.githubusercontent.com/render/math?math=c'_{i_1}X)+..+![formula](https://render.githubusercontent.com/render/math?math=c_{i_(m'-1)}X^(m'-1)) (mod q), where ![formula](https://render.githubusercontent.com/render/math?math=c_{i_0}) is participant's i secret share from time i.
3. Calculate new shares to n' and communicate via secure channels
4. Each ![formula](https://render.githubusercontent.com/render/math?math={i}\in{n'}) can calculate his groups secret share (for time i+1) such: ![formula](https://render.githubusercontent.com/render/math?math=\sum_{1}^{n'}b_i*s'_{i_j}) where ![formula](https://render.githubusercontent.com/render/math?math=b_i) is the lagrange coefficient.

## Proactive Secret Sharing (PSS) - Security model

## Proactive secret sharing libraries
* CHURP
	* [github](https://github.com/CHURPTeam/CHURP)
	* [paper](https://eprint.iacr.org/2019/017.pdf) 
* [Binance TSS](https://github.com/binance-chain/tss-lib)
* [Unbound blockchain crypto MPC](https://github.com/unbound-tech/blockchain-crypto-mpc)