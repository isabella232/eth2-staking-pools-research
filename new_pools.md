# Consensus - Creating New Pools
[![blox.io](https://s3.us-east-2.amazonaws.com/app-files.blox.io/static/media/powered_by.png)](https://blox.io)


The process of creating a new validator which can be managed by a pool goes through the same process as any other validator is created. For more details see [eth 2.0 spec](https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/validator.md#becoming-a-validator)
The special difference in creating a pool controlled validator is in the fact that the validator's keys are created via a [DKG scheme](https://github.com/bloxapp/eth2-staking-pools-research/blob/master/dkg.md) and that the creation itself could be sensitive to different attack vectors.

### Potential vulnerabilities:
- Minority Hijacking - a pool is controlled by its participants, any action requires 2/3 to sign it. If the participants set is static, it increases the likelihood of collusion. A minority of participants > 1/3 could hijack the pool itself for ransom. Although they can't "hurt" the pool, they could prevent it from doing its job and inflict penalties for the rest of the pool.
- Majorrity Hijacking - same as minority hijacking but a large amount of participants >= 2/3 can collude over time and simply distribuite the funds between themselves.
- Majority Hijacking @ creation - when creating a new pool, a set of participants is radomly grouped to create the pool. If the potential set of selection is not big enough, an attacker could easiliy create himself a majority and steal theminority's funds
- 