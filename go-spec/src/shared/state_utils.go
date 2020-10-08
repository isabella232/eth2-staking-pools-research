package shared

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared/params"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/ulule/deepcopier"
)

func CopyState(state *core.State) *core.State {
	if state == nil {
		return nil
	}

	ret := &core.State{}

	ret.CurrentSlot = state.CurrentSlot

	ret.BlockRoots = make([]*core.SlotAndBytes, len(state.BlockRoots))
	for i, r := range state.BlockRoots {
		ret.BlockRoots[i] = &core.SlotAndBytes{}
		deepcopier.Copy(r).To(ret.BlockRoots[i])
	}

	ret.StateRoots = make([]*core.SlotAndBytes, len(state.StateRoots))
	for i, r := range state.StateRoots {
		ret.StateRoots[i] = &core.SlotAndBytes{}
		deepcopier.Copy(r).To(ret.StateRoots[i])
	}

	ret.Seeds = make([]*core.SlotAndBytes, len(state.Seeds))
	for i, r := range state.Seeds {
		ret.Seeds[i] = &core.SlotAndBytes{}
		deepcopier.Copy(r).To(ret.Seeds[i])
	}

	ret.BlockProducers = make([]*core.BlockProducer, len(state.BlockProducers))
	for i, bp := range state.BlockProducers {
		ret.BlockProducers[i] = &core.BlockProducer{}
		deepcopier.Copy(bp).To(ret.BlockProducers[i])
	}

	ret.Pools = make([]*core.Pool, len(state.Pools))
	for i, p := range state.Pools {
		ret.Pools[i] = &core.Pool{}
		deepcopier.Copy(p).To(ret.Pools[i])
	}

	ret.PreviousEpochAttestations = make([]*core.PendingAttestation, len(state.PreviousEpochAttestations))
	for i, pe := range state.PreviousEpochAttestations {
		ret.PreviousEpochAttestations[i] = &core.PendingAttestation{}
		deepcopier.Copy(pe).To(ret.PreviousEpochAttestations[i])
	}

	ret.CurrentEpochAttestations = make([]*core.PendingAttestation, len(state.CurrentEpochAttestations))
	for i, pe := range state.CurrentEpochAttestations {
		ret.CurrentEpochAttestations[i] = &core.PendingAttestation{}
		deepcopier.Copy(pe).To(ret.CurrentEpochAttestations[i])
	}

	ret.JustificationBits = make(bitfield.Bitvector4, len(state.JustificationBits))
	deepcopier.Copy(state.JustificationBits).To(ret.JustificationBits)

	if state.PreviousJustifiedCheckpoint != nil {
		ret.PreviousJustifiedCheckpoint = &core.Checkpoint{}
		deepcopier.Copy(state.PreviousJustifiedCheckpoint).To(ret.PreviousJustifiedCheckpoint)
	}

	ret.CurrentJustifiedCheckpoint = &core.Checkpoint{}
	deepcopier.Copy(state.CurrentJustifiedCheckpoint).To(ret.CurrentJustifiedCheckpoint)

	if state.FinalizedCheckpoint != nil {
		ret.FinalizedCheckpoint = &core.Checkpoint{}
		deepcopier.Copy(state.FinalizedCheckpoint).To(ret.FinalizedCheckpoint)
	}

	if state.LatestBlockHeader != nil {
		ret.LatestBlockHeader = &core.PoolBlockHeader{}
		deepcopier.Copy(state.LatestBlockHeader).To(ret.LatestBlockHeader)
	}

	return ret
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
	if epoch == 0 {
		return params.ChainConfig.GenesisSeed, nil
	}

	targetSlot := epoch * params.ChainConfig.SlotsInEpoch - 1 + params.ChainConfig.SlotsInEpoch
	for _, d := range state.Seeds {
		if d.Slot == targetSlot {
			return d.Bytes, nil
		}
	}
	return []byte{}, fmt.Errorf("seed for slot %d not found", targetSlot)
}