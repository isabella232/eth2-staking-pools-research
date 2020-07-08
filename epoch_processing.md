# Consensus - Epoch Processing
[<img src="https://www.bloxstaking.com/wp-content/uploads/2020/04/Blox-Staking_logo_blue.png" width="100">](https://www.bloxstaking.com/)

Pool participants are organized around a timinig mechanism called an epoch, corresponding to eth 2.0 epoch concept. 
The reason behind it is every epoch a validator gets a set (usually just one) of duties he needs to perform in a particular block, for a pool to successfuly execute the validator's duty for a specific epoch it needs to process epochs the same way eth 2.0 does.

## Epoch Tasks
At epoch ![formula](https://render.githubusercontent.com/render/math?math=E_i), all participants of pool ![formula](https://render.githubusercontent.com/render/math?math=P_i) need to do the following tasks (as an exampele at ![formula](https://render.githubusercontent.com/render/math?math=E_i) they need to attest block ![formula](https://render.githubusercontent.com/render/math?math=B_j)):
1) Re-distribuite shares of their [group secret share](https://github.com/bloxapp/eth2-staking-pools-research/blob/master/pool_rotation.md) to pool ![formula](https://render.githubusercontent.com/render/math?math=P_i)+1 at epoch ![formula](https://render.githubusercontent.com/render/math?math=E_i)+1
2) Sign attestation ![formula](https://render.githubusercontent.com/render/math?math=A_i), partially, and broadcast signature
3) Collect broadcasted signatures, reconstruct them, and the broadcast the reconstructed signature to the beacon-chain
4) Reconstruct the participant's [group secret share](https://github.com/bloxapp/eth2-staking-pools-research/blob/master/pool_rotation.md) for ![formula](https://render.githubusercontent.com/render/math?math=E_i)+1 from shares recieved in ![formula](https://render.githubusercontent.com/render/math?math=E_i)

## Timinig Tasks
All the above tasks can't be done togteher as they have dependencies (for example you can't broadcast the signature to the beacon-chain if all pool participants didn't sign and broadcasted thier shares)
To sync the tasks in order we've divded the epoch into 3 stages, each with its own sub-tasks.

### Epoch Init stage
Timing: 1/4 down the epoch
Sub-tasks: 1

### Epoch Mid stage
Timing: 1/2 down the epoch
Sub-tasks: 2

### Epoch End stage
Timing: 3/4 down the epoch
Sub-tasks: 3 + 4