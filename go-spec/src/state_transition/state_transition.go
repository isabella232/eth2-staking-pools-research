package state_transition

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
)

type IStateTransition interface {
	// Should be called BEFORE applying the new block
	// Validates the block body and header and their contents.
	PreApplyValidateBlock(state *core.State, body *core.BlockBody) error
	// Should be called AFTER applying the new block
	PostApplyValidateBlock(newState *core.State, header *core.BlockHeader, body *core.BlockBody) error
	// Give a block, apply it's transactions on the current state.
	// Returns a new state post block.
	ApplyBlock(oldState *core.State, head *core.BlockHeader, newBlockBody *core.BlockBody) (newState *core.State, err error)

	ProcessExecutionSummaries(state *core.State, summaries []*core.ExecutionSummary) error
	ProcessNewPoolRequests(state *core.State, summaries []*core.CreateNewPoolRequest) error
}

type StateTransition struct {}
func NewStateTransition() *StateTransition { return &StateTransition{} }



// A helper function to insert the post block state root to the block body
func CalculateAndInsertStateRootToBlock(state *core.State, body *core.BlockBody) ([]byte, error) {
	st := NewStateTransition()
	newState, err := st.ApplyBlock(state, body)
	if err != nil {
		return []byte{}, err
	}

	root := core.GetStateRoot(newState, newState.CurrentEpoch)
	if len(root) == 0 {
		return []byte{}, fmt.Errorf("could not find statet root for epoch %d", newState.CurrentEpoch)
	}

	return root[:], nil
}