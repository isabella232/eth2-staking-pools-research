package block

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/prysmaticlabs/go-ssz"
)

type BlockBody struct {
	Proposer           uint64
	EpochNumber        uint64
	ExecutionSummaries []*PoolExecutionSummary
	NewPoolReq         []*CreatePoolRequest
	StateRoot          []byte
	ParentBlockRoot    []byte
}

func NewBlockBody(
	Proposer uint64,
	number uint64,
	stateRoot [32]byte,
	summary []*PoolExecutionSummary,
	newPoolReq []*CreatePoolRequest,
	parentBlockRoot []byte,
) *BlockBody {
	return &BlockBody{
		Proposer:           Proposer,
		EpochNumber:        number,
		ExecutionSummaries: summary,
		NewPoolReq:         newPoolReq,
		StateRoot:          stateRoot[:],
		ParentBlockRoot:    parentBlockRoot,
	}
}

func (body *BlockBody) GetEpochNumber() uint64 {
	return body.EpochNumber
}

func (body *BlockBody) GetProposer() uint64 {
	return body.Proposer
}

func (body *BlockBody) GetExecutionSummaries() []core.IExecutionSummary {
	ret := make([]core.IExecutionSummary, len(body.ExecutionSummaries))
	for i, d := range body.ExecutionSummaries {
		ret[i] = core.IExecutionSummary(d)
	}
	return ret
}

func (body *BlockBody) GetNewPoolRequests() []core.ICreatePoolRequest {
	ret := make([]core.ICreatePoolRequest, len(body.NewPoolReq))
	for i, d := range body.NewPoolReq {
		ret[i] = core.ICreatePoolRequest(d)
	}
	return ret
}

func (body *BlockBody) GetStateRoot() []byte {
	return body.StateRoot
}

func (body *BlockBody) GetParentBlockRoot() []byte {
	return body.ParentBlockRoot
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

