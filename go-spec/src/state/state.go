package state

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/prysmaticlabs/go-ssz"
)

type State struct {
	pools          []core.IPool
	currentEpoch   uint64
	headBlockHeader core.IBlockHeader
	blockProducers []core.IBlockProducer
	seed           [32]byte
}

func (state *State) Root() ([32]byte,error) {
	return ssz.HashTreeRoot(state)
}

func (state *State) GetPools() []core.IPool {
	return state.pools
}

func (state *State) GetPool(id uint64) core.IPool {
	for _, p := range state.pools {
		if p.GetId() == id {
			return p
		}
	}
	return nil
}

func (state *State) AddNewPool(pool core.IPool) error {
	if found := state.GetPool(pool.GetId()); found != nil {
		return fmt.Errorf("pool already exists")
	}

	state.pools = append(state.pools, pool)
	return nil
}

func (state *State) GetBlockProducers() []core.IBlockProducer {
	return state.blockProducers
}

func (state *State) GetBlockProducer(id uint64) core.IBlockProducer {
	for _, bp := range state.GetBlockProducers() {
		if bp.GetId() == id {
			return bp
		}
	}
	return nil
}

func (state *State) GetCurrentEpoch() uint64 {
	return state.currentEpoch
}

func (state *State) GetHeadBlockHeader() core.IBlockHeader {
	return state.headBlockHeader
}

func (state *State) SetHeadBlockHeader(header core.IBlockHeader){
	state.headBlockHeader = header
}

func (state *State) GetSeed() [32]byte {
	return state.seed
}

func (state *State) SetSeed(seed [32]byte) {
	state.seed = seed
}

func (state *State) GetPastSeed(epoch uint64) [32]byte {
	return [32]byte{}
}


func (state *State) Copy() (core.IState, error) {
	copiedPools := make([]core.IPool, len(state.pools))
	for i, p := range state.pools {
		newP, err := p.Copy()
		if err != nil {
			return nil, err
		}
		copiedPools[i] = newP
	}

	copiedBps := make([]core.IBlockProducer, len(state.blockProducers))
	for i, bp := range state.blockProducers {
		newBP, err := bp.Copy()
		if err != nil {
			return nil, err
		}
		copiedBps[i] = newBP
	}

	return &State{
		pools:          copiedPools,
		blockProducers: copiedBps,
		seed:           state.seed,
	}, nil
}

