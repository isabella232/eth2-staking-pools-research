package state_transition

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared/params"
)

func (st *StateTransition) ApplySlot(state *core.State, body *core.BlockBody) error {
	if isLastSlotOfEpoch(body.Slot) {
		// TODO
	}
	return nil
}

func isLastSlotOfEpoch(slot uint64) bool {
	return (slot+1) % params.ChainConfig.SlotsInEpoch == 0
}
