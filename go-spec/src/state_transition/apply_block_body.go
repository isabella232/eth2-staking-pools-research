package state_transition

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/ulule/deepcopier"
)

func (st *StateTransition) ApplyBlockBody(oldState *core.State, newBlockHeader *core.BlockHeader, newBlockBody *core.BlockBody) (newState *core.State, err error) {
	newState = &core.State{}
	deepcopier.Copy(oldState).To(newState)

	// validate
	if err := st.ValidateBlock(newState, newBlockHeader, newBlockBody); err != nil {
		return nil,err
	}

	// process
	if err := st.ProcessExecutionSummaries(newState, newBlockBody.ExecutionSummaries); err != nil {
		return nil,err
	}
	if err := st.ProcessNewPoolRequests(newState, newBlockBody.NewPoolReq); err != nil {
		return nil,err
	}

	// bump epoch
	newState.CurrentEpoch += 1
	// apply seed
	newSeed, err := shared.MixSeed(
		shared.SliceToByte32(newState.Seeds[oldState.GetCurrentEpoch()]), // previous seed
		shared.SliceToByte32(newBlockHeader.Signature[:32])) // TODO - use something else than the sig
	if err != nil {
		return nil, err
	}
	newState.Seeds[newState.CurrentEpoch] = newSeed[:]
	// add block root
	root, err := ssz.HashTreeRoot(newBlockBody)
	if err != nil {
		return nil, err
	}
	newState.BlockRoots[newState.CurrentEpoch] = root[:]
	// state root
	root, err = ssz.HashTreeRoot(newState)
	if err != nil {
		return nil, err
	}
	newState.StateRoots[newState.CurrentEpoch] = root[:]

	return newState, nil
}
