package block

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
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
	poolId uint64
	epoch  uint64 //
	duties []core.IBeaconDuty
}

func NewExecutionSummary(
	poolId uint64,
	epoch  uint64,
	duties []core.IBeaconDuty,
	) *PoolExecutionSummary {
	return &PoolExecutionSummary{
		poolId: poolId,
		epoch:  epoch,
		duties: duties,
	}
}

//func GeneratePoolSummary(
//	poolId uint64,
//	epoch uint64,
//	state *state.State,
//	) (*PoolExecutionSummary, error) {
//	// get pool and its info
//	pool := state.GetPool(poolId)
//
//	// build duties and their execution summary
//	duties, err := helperFunc.FetchExecutedDuties(pool.GetPubKey(), epoch)
//	if err != nil {
//		return nil, err
//	}
//	for _, duty := range duties {
//		duty.finalized, err = helperFunc.WasDutyIncluded(pool.GetPubKey(), epoch, duty)
//		if err != nil {
//			return nil, err
//		}
//
//		duty.participation,err = helperFunc.PoolExecutionStats(poolId, epoch, duty)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return &PoolExecutionSummary{
//		PoolId: poolId,
//		Epoch:  epoch,
//		Duties: duties,
//	}, nil
//}

func (summary *PoolExecutionSummary) GetPoolId() uint64 {
	return summary.poolId
}

func (summary *PoolExecutionSummary) GetEpoch() uint64 {
	return summary.epoch
}

func (summary *PoolExecutionSummary) GetDuties() []core.IBeaconDuty {
	return summary.duties
}

func (summary *PoolExecutionSummary) ApplyOnState(state core.IState) error {
	pool := state.GetPool(summary.GetPoolId())
	executors := pool.GetSortedExecutors()

	for _, duty := range summary.GetDuties() {
		switch duty.GetType() {
		case 0: // attestation
			for i:=0 ; i < int(core.TestConfig().PoolExecutorsNumber) ; i++ {
				bp := state.GetBlockProducer(executors[i])
				if bp == nil {
					return fmt.Errorf("BP %d not found", executors[i])
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
				bp := state.GetBlockProducer(executors[i])
				if bp == nil {
					return fmt.Errorf("BP %d not found", executors[i])
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