# Consensus - Pool Duties
A pool is simply an active eth 2.0 validator identified by a BLS12-381 public key. 
Validators get rotated between committess every epoch (epoch is 6.5 min long,  32 blocks of 12 sec), every commitee has the responsibility of advancing a speicifc eth shard and cmomitting it to the beacon-chain. 
As the eth 2.0 protocol dictates, each validator has new assignments (in his committee) every epoch for which he is rewarded or penalized for.
- Attest to a new proposed block
- Propose a new block
- Aggregate, vereify and broadcast signatures for a proposed block 

Those duties are also the duties of the staking pool in charge of a specific validator. 
To execute the validator's duty, a round leader ![formula](https://render.githubusercontent.com/render/math?math=l_{e_i}) is choosen from all the participants of pool ![formula](https://render.githubusercontent.com/render/math?math=p_{e_i}) in epoch ![formula](https://render.githubusercontent.com/render/math?math=e_i).

