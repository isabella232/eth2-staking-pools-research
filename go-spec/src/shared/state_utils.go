package shared

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
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

	ret.Randao = make([]*core.SlotAndBytes, len(state.Randao))
	for i, r := range state.Randao {
		ret.Randao[i] = &core.SlotAndBytes{}
		deepcopier.Copy(r).To(ret.Randao[i])
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

// will return nil if not found or inactive
func GetBlockProducer(state *core.State, id uint64) *core.BlockProducer {
	for _, p := range state.BlockProducers {
		if p.GetId() == id && p.Active {
			return p
		}
	}
	return nil
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

