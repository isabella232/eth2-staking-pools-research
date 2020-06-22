# Distribuited key generation (DKG) and redistribuition
[![blox.io](https://s3.us-east-2.amazonaws.com/app-files.blox.io/static/media/powered_by.png)](https://blox.io)


This section describes a standalone eth 1.0 secure DKG process for BLS12-381 keys and the math behind redistribuiting those keys to a different group as part of the participants pool rotation process.

### Basics
A DKG is a process by which a group of participants can generate a shared secret between themselves without trust. The end goal is a distribuited secret which is never constructed by any of the participants, not in the making of it nor when collectivley signing. 
Requirements:
- Correctness - the constructed key can sign, verify and is distribuited uniformly
- Secrecy - no information of the secret can be learned by an adversary
- Efficient
- Robust - can be constructed in the face of an adversary

The aggregative properties of BLS signatures enable us to use a threshold scheme with a distribuited secret. 
For the threshold scheme we will use Sahmir's scheme and for the DKG we will use Gennaro, Jarecki, Krawczyk and Rabin.

### BLS and threshold signatures
There are 2 properties of BLS signatures that make it ideal for threshold (and DKG) schemes, aggregatable and deterministic. 
We can split a single BLS key into shares via Shamir's secret sharing (SSS) scheme and a DKG scheme. The below scheme is inspired and based on [Gennaro, Jarecki, Krawczyk and Rabin's](https://link.springer.com/content/pdf/10.1007%2F3-540-48910-X_21.pdf) secure verifiable secret sharing (VSS) and [Phillip Schindler's](https://github.com/PhilippSchindler/ethdkg/blob/master/paper/ethdkg.pdf) adaptation to ETH EVM.

Definitions
- G: generator #1 for G1
- H: generator #2 for G1
- MAX_POOL_SIZE: max participants in each pool
- MIN_THRESHOLD: min amonut of participants, out of MAX_POOL_SIZE, that must sign in behalf of the pool
- QUORUM: set of participants that their shares were not disputed 

Protocol:
1. Each participant generates a randome secret, then, transfers amount S of ETH to a contract with public key (BLS12-381) of his secret.
2. Each participant generates a random polynomial of degree MIN_THRESHOLD, calculates a share for each of the other participants, SHARES[MAX_POOL_SIZE]. See below.
3. Each participant creates a commitment to the randome polynomial, see below.
4. Each participant broadcasts the shares (individually encrypted for the recipient, TBD) and his polynomial commitment.
5. Each participant can verifiy the shares he recieved (TBD)
6. Non disputed participants form the QUORUM
7. The Joint-Feldman scheme could now contruct the individually calculated shares by simply calculating (individually) the product of all recieved sahres. According to [Gennaro, Jarecki, Krawczyk and Rabin's paper](https://link.springer.com/content/pdf/10.1007%2F3-540-48910-X_21.pdf) that could result in a potential attack which results in a non uniformly distribuited secret. See below
8. To prevent said attack, all qualified parties (QUORUM) will broadcast the calculated shares but with a different generator (H) + a proof that the secret is equal betweeen g^s and h^s [see section 5.3 for moroe info](https://github.com/PhilippSchindler/ethdkg/blob/master/paper/ethdkg.pdf)
9. The master public key can then be calculated individually by each member. That public key is the pool's validator public key.

An eth EVM example could be found [here](https://github.com/PhilippSchindler/ethdkg). This example uses the BN128 curve but with the recent [EIP-2537](https://github.com/ethereum/EIPs/pull/2537) the same operations can now be done for BLS12-381

### Share Creation
Giving polynomial ![formula](https://render.githubusercontent.com/render/math?math=F(x)=c_{i_0})+![formula](https://render.githubusercontent.com/render/math?math=c_{i_1}X)+..+![formula](https://render.githubusercontent.com/render/math?math=c_{i_t}X^t) (mod q), where ![formula](https://render.githubusercontent.com/render/math?math=F(0)) is a randomly generated secret.
A participant will distribuite for other participants their shares, accoridng to their index.
![formula](https://render.githubusercontent.com/render/math?math=F(1)),![formula](https://render.githubusercontent.com/render/math?math=F(2)), ..![formula](https://render.githubusercontent.com/render/math?math=F(t))


### Polynomial Committment
Giving polynomial ![formula](https://render.githubusercontent.com/render/math?math=F(x)=c_{i_0})+![formula](https://render.githubusercontent.com/render/math?math=c_{i_1}X)+..+![formula](https://render.githubusercontent.com/render/math?math=c_{i_t}X^t) (mod q), 
a commitment to its coefficients is as followed: ![formula](https://render.githubusercontent.com/render/math?math=g^(c_{i_0})),![formula](https://render.githubusercontent.com/render/math?math=g^(c_{i_1}X)),..,![formula](https://render.githubusercontent.com/render/math?math=g^(c_{i_t}X^t))

### Joint-Feldman Key Recovery
One the QUORUM is formed, each member simply calculates his individual group secret key (gsk), ![formula](https://render.githubusercontent.com/render/math?math=gsk_i=\sum_{1}^{t}s_{i_j}), where ![formula](https://render.githubusercontent.com/render/math?math=s_{i_j}) is the individual share he recieved from participant j.
The group's master public key, ![formula](https://render.githubusercontent.com/render/math?math=mpk=\prod_{1}^{t}g^(s_i))

### Gennaro, Jarecki, Krawczyk and Rabin's DKG
[Reference implementation for BN128 curves DKG on eth](https://github.com/PhilippSchindler/ethdkg)
Each participant in QUORUM will broadcast  ![formula](https://render.githubusercontent.com/render/math?math=h^(s_i)) and a non-interactive ZK proof that for ![formula](https://render.githubusercontent.com/render/math?math=h^(s_i)) and ![formula](https://render.githubusercontent.com/render/math?math=g^(s_i)),![formula](https://render.githubusercontent.com/render/math?math=s_i) is equal. See section 5.2 [here](https://github.com/PhilippSchindler/ethdkg/blob/master/paper/ethdkg.pdf).
The master public key, ![formula](https://render.githubusercontent.com/render/math?math=mpk=\prod_{1}^{t}h^(s_i))

### Group Signature
Each member in QUORUM will partially sign a message, m. ![formula](https://render.githubusercontent.com/render/math?math=\sigma_i=H(m)^(gsk_i))
The entire group can aggregate the signatures, ![formula](https://render.githubusercontent.com/render/math?math=\sigma=\prod_{1}^{t}\sigma_i). The aggregate signature will be valid if the group threshold is achieved.

### Share Encryption
TBD

### Share Verification

TBD