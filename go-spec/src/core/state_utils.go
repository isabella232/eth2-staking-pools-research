package core

import (
	"fmt"
	"github.com/ulule/deepcopier"
)

func CopyState(state *State) *State {
	newBlockRoots := make([]*EpochAndBytes, len(state.BlockRoots))
	for i, r := range state.BlockRoots {
		newBlockRoots[i] = &EpochAndBytes{}
		deepcopier.Copy(r).To(newBlockRoots[i])
	}

	newStateRoots := make([]*EpochAndBytes, len(state.StateRoots))
	for i, r := range state.StateRoots {
		newStateRoots[i] = &EpochAndBytes{}
		deepcopier.Copy(r).To(newStateRoots[i])
	}

	newSeeds := make([]*EpochAndBytes, len(state.Seeds))
	for i, r := range state.Seeds {
		newSeeds[i] = &EpochAndBytes{}
		deepcopier.Copy(r).To(newSeeds[i])
	}

	newBPs := make([]*BlockProducer, len(state.BlockProducers))
	for i, bp := range state.BlockProducers {
		newBPs[i] = &BlockProducer{}
		deepcopier.Copy(bp).To(newBPs[i])
	}

	newPools := make([]*Pool, len(state.Pools))
	for i, p := range state.Pools {
		newPools[i] = &Pool{}
		deepcopier.Copy(p).To(newPools[i])
	}

	return &State{
		GenesisTime:    state.GenesisTime,
		CurrentEpoch:   state.CurrentEpoch,
		BlockRoots:     newBlockRoots,
		StateRoots:     newStateRoots,
		Seeds:          newSeeds,
		BlockProducers: newBPs,
		Pools:          newPools,
	}
}

// will return an 0 length byte array if not found
func GetStateRoot(state *State, epoch uint64) []byte {
	for _, r := range state.StateRoots {
		if r.Epoch == epoch {
			return r.Bytes
		}
	}
	return []byte{}
}

func DecreaseBPBalance(bp *BlockProducer, change uint64) error {
	if bp.CDTBalance < change {
		return fmt.Errorf("BP %d dosen't have enonugh Balance (%d) to decrease (%d)", bp.Id, bp.CDTBalance, change)
	}

	bp.CDTBalance -= change
	return nil
}

func IncreaseBPBalance(bp *BlockProducer, change uint64) {
	bp.CDTBalance += change
}

// will return nil if not found or inactive
func GetBlockProducer(state *State, id uint64) *BlockProducer {
	for _, p := range state.BlockProducers {
		if p.GetId() == id && p.Active {
			return p
		}
	}
	return nil
}

func GetActiveBlockProducers(state *State, epoch uint64) []uint64 {
	var activeBps []uint64
	for _, bp := range state.BlockProducers {
		if bp.Active || bp.GetExitEpoch() > epoch {
			activeBps = append(activeBps, bp.GetId())
		}
	}
	return activeBps
}

// will return nil if not found
func GetPool(state *State, id uint64) *Pool {
	for _, p := range state.Pools {
		if p.GetId() == id {
			return p
		}
	}
	return nil
}

// will return error if not found
func GetSeed(state *State, epoch uint64) ([]byte, error) {
	for _, d := range state.Seeds {
		if d.Epoch == epoch {
			return d.Bytes, nil
		}
	}
	return []byte{}, fmt.Errorf("seed for epoch %d not found", epoch)
}

