package src

func GenerateRandomState() *State {
	pools := make([]*Pool, 5)

	//
	bps := make([]*BlockProducer, len(pools) * int(TestConfig().PoolExecutorsNumber))
	for i := 0 ; i < len(bps) ; i++ {
		bps[i] = &BlockProducer{
			Id:      uint64(i),
			Balance: 1000,
			Stake:   0,
			Slashed: false,
		}
	}

	//
	for i := 0 ; i < len(pools) ; i++ {
		executors := make([]uint64, TestConfig().PoolExecutorsNumber)
		for j := 0 ; j < int(TestConfig().PoolExecutorsNumber) ; j++ {
			executors[j] = bps[i*int(TestConfig().PoolExecutorsNumber) + j].Id
		} // no need to sort as they are already

		pools[i] = &Pool{
			Id:              uint64(i),
			SortedExecutors: executors,
		}
	}

	return &State{
		Pools:           pools,
		BlockRoots:      nil,
		HeadBlockHeader: nil,
		BlockProducers:  bps,
		Seed:            []byte("seed"),
	}
}
