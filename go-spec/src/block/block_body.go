package block

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/state"
	"github.com/prysmaticlabs/go-ssz"
)

type BlockBody struct {
	proposer           uint64
	epochNumber        uint64
	executionSummaries []core.IExecutionSummary
	newPoolReq         []core.ICreatePoolRequest
	stateRoot          []byte
	parentBlockRoot    []byte
}

func NewBlockBody(
	Proposer uint64,
	number uint64,
	state *state.State,
	summary []core.IExecutionSummary,
	parentBlockRoot []byte,
) (*BlockBody, error) {
	stateRoot,err := state.Root()
	if err != nil {
		return nil, err
	}

	return &BlockBody{
		proposer:           Proposer,
		epochNumber:        number,
		executionSummaries: summary,
		stateRoot:          stateRoot[:],
		parentBlockRoot:    parentBlockRoot,
	}, nil
}

func (body *BlockBody) GetEpochNumber() uint64 {
	return body.epochNumber
}

func (body *BlockBody) GetProposer() uint64 {
	return body.proposer
}

func (body *BlockBody) GetExecutionSummaries() []core.IExecutionSummary {
	return body.executionSummaries
}

func (body *BlockBody) GetNewPoolRequests() []core.ICreatePoolRequest {
	return body.newPoolReq
}

func (body *BlockBody) GetStateRoot() []byte {
	return body.stateRoot
}

func (body *BlockBody) GetParentBlockRoot() []byte {
	return body.parentBlockRoot
}

func (body *BlockBody) Root() ([]byte, error) {
	ret, err := ssz.HashTreeRoot(body)
	if err != nil {
		return nil, err
	}
	return ret[:],nil
}

func (body *BlockBody) Validate() error {
	return nil
}

