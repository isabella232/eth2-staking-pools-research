package state_transition

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
)

type IStateTransition interface {
	ValidateBlock(state *core.State, header *core.BlockHeader, body *core.BlockBody) error
	ApplyBlockBody(oldState *core.State, newBlockHeader *core.BlockHeader, newBlockBody *core.BlockBody) (newState *core.State, err error)

	ProcessExecutionSummaries(state *core.State, summaries []*core.ExecutionSummary) error
	ProcessNewPoolRequests(state *core.State, summaries []*core.CreateNewPoolRequest) error
}

type StateTransition struct {

}

