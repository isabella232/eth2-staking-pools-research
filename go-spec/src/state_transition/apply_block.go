package state_transition

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/prysmaticlabs/go-ssz"
)

func (st *StateTransition) ApplyBlock(oldState *core.State, body *core.BlockBody) (newState *core.State, err error) {
	newState = shared.CopyState(oldState)

	// bump epoch
	previousSlot := newState.CurrentSlot
	newState.CurrentSlot = body.Slot

	// process
	if err := st.ProcessNewPoolRequests(newState, body.NewPoolReq); err != nil {
		return nil,err
	}
	if err := st.ProcessBlockAttestations(newState, body); err != nil {
		return nil,err
	}

	// apply seed
	prevSeed, err := shared.GetSeed(newState, previousSlot)
	if err != nil {
		return nil, err
	}
	newSeed, err := shared.MixSeed(
		shared.SliceToByte32(prevSeed), // previous seed
		shared.SliceToByte32(body.Randao[:32]))
	if err != nil {
		return nil, err
	}
	newState.Seeds = append(newState.Seeds, &core.SlotAndBytes{
		Slot:                newState.CurrentSlot,
		Bytes:               newSeed[:],
	})
	// add block root
	root, err := ssz.HashTreeRoot(body)
	if err != nil {
		return nil, err
	}
	newState.BlockRoots = append(newState.BlockRoots, &core.SlotAndBytes{
		Slot:                newState.CurrentSlot,
		Bytes:               root[:],
	})
	// state root
	root, err = ssz.HashTreeRoot(newState)
	if err != nil {
		return nil, err
	}
	newState.StateRoots = append(newState.StateRoots, &core.SlotAndBytes{
		Slot:                 newState.CurrentSlot,
		Bytes:                root[:],
	})

	if err := st.ApplySlot(newState, body); err != nil {
		return nil, err
	}

	return newState, nil
}
