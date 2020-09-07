package src

type PoolExecutionSummary struct {
	PoolId        	[]byte // eth2 validator pubkey
	StartingEpoch 	uint64 // a.k.a previous epoch
	EndEpoch      	uint64 //
	Performance   	map[*BeaconDuty][16]byte // for every duty specify an array of 128 bits (16 bytes) of who participated in the execution of that duty
}

// will calculate rewards/ penalties and apply them onto the state
func (summary *PoolExecutionSummary) ApplyOnState(state *State) error {
	pool, err := GetPool(state, summary.PoolId)
	if err != nil {
		return err
	}

	for duty, whoExecuted := range summary.Performance {
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
					if IsBitSet(whoExecuted[:], uint64(i)) {
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
					if IsBitSet(whoExecuted[:], uint64(i)) {
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