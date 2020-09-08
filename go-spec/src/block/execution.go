package block

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/state"
)

/**
	This object is crucial for the honest operation of executors assigned to a pool.
	It does 2 things:
	1) checks which duties the pool had in a specific epoch
	2) submits an array of 16 bytes (128 bits) which represents a 1 for each executor that participated in executing the duty
	   and 0 for each that didn't.

	This helps assigning rewards/ penalties for executors.
 */
type PoolExecutionSummary struct {
	PoolId uint64
	Epoch  uint64 //
	Duties []core.IBeaconDuty
}

func GeneratePoolSummary(
	poolId uint64,
	epoch uint64,
	state *state.State,
	helperFunc src.NonSpecFunctions,
	) (*PoolExecutionSummary, error) {
	// get pool and its info
	pool := state.GetPool(poolId)

	// build duties and their execution summary
	duties, err := helperFunc.FetchExecutedDuties(pool.GetPubKey(), epoch)
	if err != nil {
		return nil, err
	}
	for _, duty := range duties {
		duty.finalized, err = helperFunc.WasDutyIncluded(pool.GetPubKey(), epoch, duty)
		if err != nil {
			return nil, err
		}

		duty.participation,err = helperFunc.PoolExecutionStats(poolId, epoch, duty)
		if err != nil {
			return nil, err
		}
	}

	return &PoolExecutionSummary{
		PoolId: poolId,
		Epoch:  epoch,
		Duties: duties,
	}, nil
}

func (summary *PoolExecutionSummary) GetPoolId() uint64 {
	return summary.PoolId
}

func (summary *PoolExecutionSummary) GetEpoch() uint64 {
	return summary.Epoch
}

func (summary *PoolExecutionSummary) GetDuties() []core.IBeaconDuty {
	return summary.Duties
}

func (summary *PoolExecutionSummary) ApplyOnState(state core.IState) error {
	pool := state.GetPool(summary.PoolId)

	for _, duty := range summary.Duties {
		switch duty.GetType() {
		case 0: // attestation
			for i:=0 ; i < int(core.TestConfig().PoolExecutorsNumber) ; i++ {
				executor := pool.GetSortedExecutors()[i]
				bp := state.GetBlockProducer(executor)
				if bp == nil {
					return fmt.Errorf("BP %d not found", executor)
				}

				if !duty.IsFinalized() {
					_,err := bp.DecreaseBalance(2*core.TestConfig().BaseEth2DutyReward)
					if err != nil {
						return err
					}
				} else {
					participation := duty.GetParticipation()
					if shared.IsBitSet(participation[:], uint64(i)) {
						_,err := bp.IncreaseBalance(core.TestConfig().BaseEth2DutyReward)
						if err != nil {
							return err
						}
					} else {
						_,err := bp.DecreaseBalance(core.TestConfig().BaseEth2DutyReward)
						if err != nil {
							return err
						}
					}
				}
			}
		case 1: // proposal
			for i:=0 ; i < int(core.TestConfig().PoolExecutorsNumber) ; i++ {
				executor := pool.GetSortedExecutors()[i]
				bp := state.GetBlockProducer(executor)
				if bp == nil {
					return fmt.Errorf("BP %d not found", executor)
				}

				if !duty.IsFinalized() {
					_,err := bp.DecreaseBalance(4*core.TestConfig().BaseEth2DutyReward)
					if err != nil {
						return err
					}
				} else {
					participation := duty.GetParticipation()
					if shared.IsBitSet(participation[:], uint64(i)) {
						_,err := bp.IncreaseBalance(2*core.TestConfig().BaseEth2DutyReward)
						if err != nil {
							return err
						}
					} else {
						_,err := bp.DecreaseBalance(2*core.TestConfig().BaseEth2DutyReward)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}