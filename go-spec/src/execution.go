package src

import "github.com/bloxapp/eth2-staking-pools-research/go-spec/src/state"

type BeaconDuty struct {
	Type				uint8 // 0 - attestation, 1 - block proposal
	Committee 			uint64
	Slot				uint64
	Included			bool // whether or not it was included in the beacon chain (the pool earned reward from it)
	Executors			[16]byte // 128 bit of the executors (by order) which executed this duty
}

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
	Duties []*BeaconDuty
}

func GeneratePoolSummary(
	poolId uint64,
	epoch uint64,
	state *state.State,
	helperFunc NonSpecFunctions,
	) (*PoolExecutionSummary, error) {
	// get pool and its info
	pool := state.GetPool(poolId)

	// build duties and their execution summary
	duties, err := helperFunc.FetchExecutedDuties(pool.PubKey, epoch)
	if err != nil {
		return nil, err
	}
	for _, duty := range duties {
		duty.Included, err = helperFunc.WasDutyIncluded(pool.PubKey, epoch, duty)
		if err != nil {
			return nil, err
		}

		duty.Executors,err = helperFunc.PoolExecutionStats(poolId, epoch, duty)
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

// will calculate rewards/ penalties and apply them onto the state
func (summary *PoolExecutionSummary) ApplyOnState(state *state.State) error {
	pool := state.GetPool(summary.PoolId)

	for _, duty := range summary.Duties {
		switch duty.Type {
		case 0: // attestation
			for i:=0 ; i < int(TestConfig().PoolExecutorsNumber) ; i++ {
				executor := pool.SortedExecutors[i]
				if !duty.Included {
					_,err := state.DecreaseBlockProducerBalance(executor, 2*TestConfig().BaseEth2DutyReward)
					if err != nil {
						return err
					}
				} else {
					if IsBitSet(duty.Executors[:], uint64(i)) {
						_,err := state.IncreaseBlockProducerBalance(executor, TestConfig().BaseEth2DutyReward)
						if err != nil {
							return err
						}
					} else {
						_,err := state.DecreaseBlockProducerBalance(executor, TestConfig().BaseEth2DutyReward)
						if err != nil {
							return err
						}
					}
				}
			}
		case 1: // proposal
			for i:=0 ; i < int(TestConfig().PoolExecutorsNumber) ; i++ {
				executor := pool.SortedExecutors[i]
				if !duty.Included {
					_,err := state.DecreaseBlockProducerBalance(executor, 4*TestConfig().BaseEth2DutyReward)
					if err != nil {
						return err
					}
				} else {
					if IsBitSet(duty.Executors[:], uint64(i)) {
						_,err := state.IncreaseBlockProducerBalance(executor, 2*TestConfig().BaseEth2DutyReward)
						if err != nil {
							return err
						}
					} else {
						_,err := state.DecreaseBlockProducerBalance(executor, 2*TestConfig().BaseEth2DutyReward)
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