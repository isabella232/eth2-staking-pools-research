# Consensus - Pool Duties
[![blox.io](https://s3.us-east-2.amazonaws.com/app-files.blox.io/static/media/powered_by.png)](https://www.bloxstaking.com)


A pool is simply an active eth 2.0 validator identified by a BLS12-381 public key. 
Validators get rotated between committess every epoch (epoch is 6.5 min long,  32 blocks of 12 sec), every commitee has the responsibility of advancing a speicifc eth shard and committing it to the beacon-chain. 
As the eth 2.0 protocol dictates, each validator has new assignments (in his committee) every epoch for which he is rewarded or penalized for.
- Attest to a new proposed block
- Propose a new block
- Aggregate, vereify and broadcast signatures for a proposed block 

Those duties are also the duties of the staking pool in charge of a specific validator. 
To execute the validator's duty, a round leader ![formula](https://render.githubusercontent.com/render/math?math=l_{e_i}) is choosen randomly from all the participants of pool ![formula](https://render.githubusercontent.com/render/math?math=p_{e_i}) in epoch ![formula](https://render.githubusercontent.com/render/math?math=e_i).
Choosing a randome ![formula](https://render.githubusercontent.com/render/math?math=l_{e_i}) could be done using the beacon-chain's RANDAO seed.

### Round leader responsibility
![formula](https://render.githubusercontent.com/render/math?math=l_{e_i}) will broadcast to the pool's participants a proposal for them to sign. The proposal depends on the duty at hand.
* Block attestation - ![formula](https://render.githubusercontent.com/render/math?math=l_{e_i}) will send an [attestation_data](https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/beacon-chain.md#attestationdata) object
* Block proposal - ![formula](https://render.githubusercontent.com/render/math?math=l_{e_i}) will send a [SignedBeaconBlock](https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/beacon-chain.md#signedbeaconblock)

If ![formula](https://render.githubusercontent.com/render/math?math=l_{e_i}) proposes a slahsable proposal, he will be slashed from the protocol.\
If ![formula](https://render.githubusercontent.com/render/math?math=l_{e_i}) fails to send a proposal in the correct time, he will be penalized solely.\
If ![formula](https://render.githubusercontent.com/render/math?math=l_{e_i}) does send a proposal but it fails to gather the threshold amount of signatures, the lost duties penalty will be divided between all participants of ![formula](https://render.githubusercontent.com/render/math?math=p_{e_i})\
If ![formula](https://render.githubusercontent.com/render/math?math=l_{e_i}) successfuly proposes and that proposal successfuly signed and broadcasted, he will recieve a special share of the reward.\
**For details see rewards and penalties section**

The round leader also has the eresponsibility of broadccasting to the rest of the protocol a summary of which participants fulfilled their duties.

### Participant's responsibility
A participant of ![formula](https://render.githubusercontent.com/render/math?math=p_{e_i}) has the responsibility of recieving a proposal from ![formula](https://render.githubusercontent.com/render/math?math=l_{e_i}), verifying it and signing it.
If ![formula](https://render.githubusercontent.com/render/math?math=p_{e_i}) signs a slashable proposal, he will be slashed from the protocol.
If ![formula](https://render.githubusercontent.com/render/math?math=p_{e_i}) does not sign a proposal he will be penalized
If ![formula](https://render.githubusercontent.com/render/math?math=p_{e_i}) signs a valid proposal, that proposal was broadcasted and the validator recieved a reward, ![formula](https://render.githubusercontent.com/render/math?math=p_{e_i}) will be rewarded proportionally.