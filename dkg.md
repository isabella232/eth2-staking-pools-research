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
There are 2 properties of BLS signatures that make it ideal for threshold (and DKG) schemes. They are aggregatable and deterministic. 
Shamir's scheme consists of creating shares 