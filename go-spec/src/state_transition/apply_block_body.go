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
	newState.CurrentEpoch = newBlockBody.Epoch
	// apply seed
	newSeed, err := shared.MixSeed(
		shared.SliceToByte32(oldState.Seeds[len(oldState.Seeds) - 1].Bytes), // previous seed
		shared.SliceToByte32(newBlockHeader.Signature[:32])) // TODO - use something else than the sig
	if err != nil {
		return nil, err
	}
	newState.Seeds = append(newState.Seeds, &core.EpochAndBytes{
		Epoch:                newState.CurrentEpoch,
		Bytes:                newSeed[:],
	})
	// add block root
	root, err := ssz.HashTreeRoot(newBlockBody)
	if err != nil {
		return nil, err
	}
	newState.BlockRoots = append(newState.BlockRoots, &core.EpochAndBytes{
		Epoch:                newState.CurrentEpoch,
		Bytes:                root[:],
	})
	// state root
	root, err = ssz.HashTreeRoot(newState)
	if err != nil {
		return nil, err
	}
	newState.StateRoots = append(newState.StateRoots, &core.EpochAndBytes{
		Epoch:                newState.CurrentEpoch,
		Bytes:                root[:],
	})

	return newState, nil
}
