package state

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/prysmaticlabs/go-ssz"
)

type State struct {
	GenesisTime    uint64
	CurrentEpoch   uint64
	BlockRoots     [][32]byte // fork choice block roots
	StateRoots     [][32]byte // fork choice state roots
	Seeds          [][32]byte
	BlockProducers []*BlockProducer
	Pools          []*Pool
	Slashings      [][]uint64 // fork choice Slashings
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
		GenesisTime:  genesisTime,
		Pools:        pools,
		CurrentEpoch: currentEpoch,
		//headBlockHeader: headBlockHeader,
		BlockProducers: blockProducers,
		Seeds:          [][32]byte{epochZeroSeed},
		BlockRoots:     [][32]byte{},
		StateRoots:     [][32]byte{},
		Slashings:      [][]uint64{},
	}
}

func (state *State) Root() ([32]byte,error) {
	return ssz.HashTreeRoot(state)
}

func (state *State) GetPools() []core.IPool {
	ret := make([]core.IPool, len(state.Pools))
	for i, d := range state.Pools {
		ret[i] = core.IPool(d)
	}
	return ret
}

func (state *State) GetPool(id uint64) core.IPool {
	for _, p := range state.Pools {
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

	state.Pools = append(state.Pools, pool)
	return nil
}

func (state *State) GetBlockProducers() []core.IBlockProducer {
	ret := make([]core.IBlockProducer, len(state.BlockProducers))
	for i, d := range state.BlockProducers {
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
	return state.CurrentEpoch
}

func (state *State) SetCurrentEpoch(epoch uint64) {
	state.CurrentEpoch = epoch
}

func (state *State) GetSeed(epoch uint64) [32]byte {
	return state.Seeds[epoch]
}

func (state *State) SetSeed(seed [32]byte, epoch uint64) {
	if uint64(len(state.Seeds)) <= epoch {
		state.Seeds = append(state.Seeds, seed)
	} else {
		state.Seeds[epoch] = seed
	}
}

func (state *State) GetPastSeed(epoch uint64) [32]byte {
	return [32]byte{}
}

func (state *State) GetBlockRoot(epoch uint64) [32]byte {
	return state.BlockRoots[epoch]
}

func (state *State) SetBlockRoot(root [32]byte, epoch uint64) {
	if uint64(len(state.BlockRoots)) <= epoch {
		state.BlockRoots = append(state.BlockRoots, root)
	} else {
		state.BlockRoots[epoch] = root
	}
}

func (state *State) GetStateRoot(epoch uint64) [32]byte {
	return state.StateRoots[epoch]
}

func (state *State) SetStateRoot(root [32]byte, epoch uint64) {
	if uint64(len(state.StateRoots)) <= epoch {
		state.StateRoots = append(state.StateRoots, root)
	} else {
		state.StateRoots[epoch] = root
	}
}

func (state *State) Copy() (core.IState, error) {
	copiedPools := make([]*Pool, len(state.Pools))
	for i, p := range state.Pools {
		copiedPools[i] = &Pool{
			Id:              p.GetId(),
			Active:          p.IsActive(),
			PubKey:          p.PubKey,
			SortedExecutors: p.GetSortedExecutors(),
		}
	}

	copiedBps := make([]*BlockProducer, len(state.BlockProducers))
	for i, bp := range state.BlockProducers {
		copiedBps[i] = &BlockProducer{
			Id:      bp.Id,
			PubKey:  bp.PubKey,
			Balance: bp.Balance,
			Stake:   bp.Stake,
			Slashed: bp.Slashed,
			Active:  bp.Active,
		}
	}

	return &State{
		Pools:          copiedPools,
		BlockProducers: copiedBps,
		Seeds:          state.Seeds,
	}, nil
}

