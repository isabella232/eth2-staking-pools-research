package state

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/prysmaticlabs/go-ssz"
)

type State struct {
	genesisTime     uint64
	currentEpoch    uint64
	blockRoots      [][32]byte // fork choice block roots
	stateRoots      [][32]byte // fork choice state roots
	seeds           [][32]byte
	blockProducers  []*BlockProducer
	pools           []*Pool
	slashings       [][]uint64 // fork choice slashings
}

func NewState(
	genesisTime uint64,
	pools []*Pool,
	currentEpoch   uint64,
	//headBlockHeader core.IBlockHeader,
	blockProducers []*BlockProducer,
	epochZeroSeed [32]byte,
	) *State {
	return &State{
		genesisTime:	 genesisTime,
		pools:           pools,
		currentEpoch:    currentEpoch,
		//headBlockHeader: headBlockHeader,
		blockProducers:  blockProducers,
		seeds:           [][32]byte{epochZeroSeed},
		blockRoots: 	 [][32]byte{},
		stateRoots: 	 [][32]byte{},
		slashings: 		 [][]uint64{},
	}
}

func (state *State) Root() ([32]byte,error) {
	// TODO - state root
	return ssz.HashTreeRoot("state root ssz")
}

func (state *State) GetPools() []core.IPool {
	ret := make([]core.IPool, len(state.pools))
	for i, d := range state.pools {
		ret[i] = core.IPool(d)
	}
	return ret
}

func (state *State) GetPool(id uint64) core.IPool {
	for _, p := range state.pools {
		if p.GetId() == id {
			return p
		}
	}
	return nil
}

func (state *State) AddNewPool(pool *Pool) error {
	if found := state.GetPool(pool.GetId()); found != nil {
		return fmt.Errorf("pool already exists")
	}

	state.pools = append(state.pools, pool)
	return nil
}

func (state *State) GetBlockProducers() []core.IBlockProducer {
	ret := make([]core.IBlockProducer, len(state.blockProducers))
	for i, d := range state.blockProducers {
		ret[i] = core.IBlockProducer(d)
	}
	return ret
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

func (state *State) SetCurrentEpoch(epoch uint64) {
	state.currentEpoch = epoch
}

func (state *State) GetSeed(epoch uint64) [32]byte {
	return state.seeds[epoch]
}

func (state *State) SetSeed(seed [32]byte, epoch uint64) {
	if uint64(len(state.seeds)) <= epoch {
		state.seeds = append(state.seeds, seed)
	} else {
		state.seeds[epoch] = seed
	}
}

func (state *State) GetPastSeed(epoch uint64) [32]byte {
	return [32]byte{}
}

func (state *State) GetBlockRoot(epoch uint64) [32]byte {
	return state.blockRoots[epoch]
}

func (state *State) SetBlockRoot(root [32]byte, epoch uint64) {
	if uint64(len(state.blockRoots)) <= epoch {
		state.blockRoots = append(state.blockRoots, root)
	} else {
		state.blockRoots[epoch] = root
	}
}

func (state *State) GetStateRoot(epoch uint64) [32]byte {
	return state.stateRoots[epoch]
}

func (state *State) SetStateRoot(root [32]byte, epoch uint64) {
	if uint64(len(state.stateRoots)) <= epoch {
		state.stateRoots = append(state.stateRoots, root)
	} else {
		state.stateRoots[epoch] = root
	}
}

func (state *State) Copy() (core.IState, error) {
	copiedPools := make([]*Pool, len(state.pools))
	for i, p := range state.pools {
		copiedPools[i] = NewPool(
				p.GetId(),
				p.IsActive(),
				p.pubKey,
				p.GetSortedExecutors(),
			)
	}

	copiedBps := make([]*BlockProducer, len(state.blockProducers))
	for i, bp := range state.blockProducers {
		copiedBps[i] = &BlockProducer{
			id:      bp.id,
			pubKey:  bp.pubKey,
			balance: bp.balance,
			stake:   bp.stake,
			slashed: bp.slashed,
			active:  bp.active,
		}
	}

	return &State{
		pools:          copiedPools,
		blockProducers: copiedBps,
		seeds:           state.seeds,
	}, nil
}

