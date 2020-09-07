package src

type PoolExecutionSummary struct {
	PoolId        	uint64
	StartingEpoch 	uint64 // a.k.a previous epoch
	EndEpoch      	uint64 //
	Duties   		[]*BeaconDuty
}

// will calculate rewards/ penalties and apply them onto the state
func (summary *PoolExecutionSummary) ApplyOnState(state *State) error {
	pool, err := GetPool(state, summary.PoolId)
	if err != nil {
		return err
	}

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
						_,err = state.DecreaseBlockProducerBalance(executor, TestConfig().BaseEth2DutyReward)
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
					_,err = state.DecreaseBlockProducerBalance(executor, 4*TestConfig().BaseEth2DutyReward)
					if err != nil {
						return err
					}
				} else {
					if IsBitSet(duty.Executors[:], uint64(i)) {
						_,err = state.IncreaseBlockProducerBalance(executor, 2*TestConfig().BaseEth2DutyReward)
						if err != nil {
							return err
						}
					} else {
						_,err = state.DecreaseBlockProducerBalance(executor, 2*TestConfig().BaseEth2DutyReward)
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