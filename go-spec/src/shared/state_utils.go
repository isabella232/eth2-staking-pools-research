package shared

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared/params"
	"github.com/ulule/deepcopier"
)

func CopyState(state *core.State) *core.State {
	newBlockRoots := make([]*core.SlotAndBytes, len(state.BlockRoots))
	for i, r := range state.BlockRoots {
		newBlockRoots[i] = &core.SlotAndBytes{}
		deepcopier.Copy(r).To(newBlockRoots[i])
	}

	newStateRoots := make([]*core.SlotAndBytes, len(state.StateRoots))
	for i, r := range state.StateRoots {
		newStateRoots[i] = &core.SlotAndBytes{}
		deepcopier.Copy(r).To(newStateRoots[i])
	}

	newSeeds := make([]*core.SlotAndBytes, len(state.Seeds))
	for i, r := range state.Seeds {
		newSeeds[i] = &core.SlotAndBytes{}
		deepcopier.Copy(r).To(newSeeds[i])
	}

	newBPs := make([]*core.BlockProducer, len(state.BlockProducers))
	for i, bp := range state.BlockProducers {
		newBPs[i] = &core.BlockProducer{}
		deepcopier.Copy(bp).To(newBPs[i])
	}

	newPools := make([]*core.Pool, len(state.Pools))
	for i, p := range state.Pools {
		newPools[i] = &core.Pool{}
		deepcopier.Copy(p).To(newPools[i])
	}

	return &core.State{
		GenesisTime:    state.GenesisTime,
		CurrentSlot:   state.CurrentSlot,
		BlockRoots:     newBlockRoots,
		StateRoots:     newStateRoots,
		Seeds:          newSeeds,
		BlockProducers: newBPs,
		Pools:          newPools,
	}
}

func GetCurrentEpoch(state *core.State) uint64 {
	return params.SlotToEpoch(state.CurrentSlot)
}

func GetPreviousEpoch(state *core.State) (uint64, error) {
	current := params.SlotToEpoch(state.CurrentSlot)
	if current != 0 {
		return current - 1, nil
	} else {
		return 0, fmt.Errorf("current peoch is 0, no previous epoch")
	}
}

// will return an 0 length byte array if not found
func GetStateRoot(state *core.State, slot uint64) []byte {
	for _, r := range state.StateRoots {
		if r.Slot == slot {
			return r.Bytes
		}
	}
	return []byte{}
}

func DecreaseBPBalance(bp *core.BlockProducer, change uint64) error {
	if bp.CDTBalance < change {
		return fmt.Errorf("BP %d dosen't have enonugh Balance (%d) to decrease (%d)", bp.Id, bp.CDTBalance, change)
	}

	bp.CDTBalance -= change
	return nil
}

func IncreaseBPBalance(bp *core.BlockProducer, change uint64) {
	bp.CDTBalance += change
}

// will return nil if not found or inactive
func GetBlockProducer(state *core.State, id uint64) *core.BlockProducer {
	for _, p := range state.BlockProducers {
		if p.GetId() == id && p.Active {
			return p
		}
	}
	return nil
}

func GetActiveBlockProducers(state *core.State, epoch uint64) []uint64 {
	var activeBps []uint64
	for _, bp := range state.BlockProducers {
		if bp.Active || bp.GetExitEpoch() > epoch {
			activeBps = append(activeBps, bp.GetId())
		}
	}
	return activeBps
}

// will return nil if not found
func GetPool(state *core.State, id uint64) *core.Pool {
	for _, p := range state.Pools {
		if p.GetId() == id {
			return p
		}
	}
	return nil
}

// Returns the seed after randao was applied on the last slot of the epoch
// will return error if not found
func GetEpochSeed(state *core.State, epoch uint64) ([]byte, error) {
	targetSlot := epoch * params.ChainConfig.SlotsInEpoch - 1 + params.ChainConfig.SlotsInEpoch
	seed, err := GetSeed(state, targetSlot)
	if err != nil {
		return []byte{}, fmt.Errorf("seed for epoch %d not found", epoch)
	}
	return seed, nil
}

// returns seed for a slot
func GetSeed(state *core.State, slot uint64) ([]byte, error) {
	for _, d := range state.Seeds {
		if d.Slot == slot {
			return d.Bytes, nil
		}
	}
	return []byte{}, fmt.Errorf("seed for slot %d not found", slot)
}

