# Distribuited key generation (DKG) and redistribuition
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

$x_{n}$