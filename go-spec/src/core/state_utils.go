package core

import (
	"fmt"
	"github.com/ulule/deepcopier"
)

func CopyState(state *State) *State {
	newBlockRoots := make([]*SlotAndBytes, len(state.BlockRoots))
	for i, r := range state.BlockRoots {
		newBlockRoots[i] = &SlotAndBytes{}
		deepcopier.Copy(r).To(newBlockRoots[i])
	}

	newStateRoots := make([]*SlotAndBytes, len(state.StateRoots))
	for i, r := range state.StateRoots {
		newStateRoots[i] = &SlotAndBytes{}
		deepcopier.Copy(r).To(newStateRoots[i])
	}

	newSeeds := make([]*SlotAndBytes, len(state.Seeds))
	for i, r := range state.Seeds {
		newSeeds[i] = &SlotAndBytes{}
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
		CurrentSlot:   state.CurrentSlot,
		BlockRoots:     newBlockRoots,
		StateRoots:     newStateRoots,
		Seeds:          newSeeds,
		BlockProducers: newBPs,
		Pools:          newPools,
	}
}

// will return an 0 length byte array if not found
func GetStateRoot(state *State, slot uint64) []byte {
	for _, r := range state.StateRoots {
		if r.Slot == slot {
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

// Returns the seed after randao was applied on the last slot of the epoch
// will return error if not found
func GetEpochSeed(state *State, epoch uint64) ([]byte, error) {
	targetSlot := epoch * TestConfig().SlotsInEpoch - 1 + TestConfig().SlotsInEpoch
	seed, err := GetSeed(state, targetSlot)
	if err != nil {
		return []byte{}, fmt.Errorf("seed for epoch %d not found", epoch)
	}
	return seed, nil
}

// returns seed for a slot
func GetSeed(state *State, slot uint64) ([]byte, error) {
	for _, d := range state.Seeds {
		if d.Slot == slot {
			return d.Bytes, nil
		}
	}
	return []byte{}, fmt.Errorf("seed for slot %d not found", slot)
}

