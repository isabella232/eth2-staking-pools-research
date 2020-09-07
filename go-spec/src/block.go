package src

import (
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-ssz"
)

//// A request struct for creating new pool credentials
//// will trigger random selection of 128 executors to DKG new pool credentials and wait for deposit
//type CreatePoolRequest struct {
//
//}
//
//// A request struct for BPs to withdraw their rewards to eth mainnet CDT contract
//type WithdrawRequest struct {
//
//}
//
//// A request to liquidate a pool, should go through validation
//type LiquidatePoolRequest struct {
//
//}
//
//// Represents slashed block producers
//type Slashing struct {
//
//}

type BlockBody struct {
	Proposer 				uint64
	Number					uint64
	PoolsExecutionSummary 	[]*PoolExecutionSummary
	//NewPoolReq			[]*CreatePoolRequest
	//WithdrawReq			[]*WithdrawRequest
	//LiquidationReq		[]*LiquidatePoolRequest
	//Slashing			[]*Slashing
	StateRoot				[]byte
	ParentBlockRoot			[]byte
}

func BuildBlockBody(
	Proposer uint64,
	number uint64,
	state *State,
	summary []*PoolExecutionSummary,
	helperFunc NonSpecFunctions,
	) (*BlockBody, error) {
	stateRoot,err := state.Root()
	if err != nil {
		return nil, err
	}

	headBlockBody, err := helperFunc.GetBlockBody(state.HeadBlockHeader.BlockRoot)
	if err != nil {
		return nil, err
	}
	parentRoot,err := headBlockBody.Root()
	if err != nil {
		return nil, err
	}

	return &BlockBody{
		Proposer:              Proposer,
		Number:                number,
		PoolsExecutionSummary: summary,
		StateRoot:             stateRoot[:],
		ParentBlockRoot:       parentRoot,
	}, nil
}

func (body *BlockBody) Validate() error {
	return nil
}

func (body *BlockBody) Root() ([]byte,error) {
	ret, err := ssz.HashTreeRoot(body)
	if err != nil {
		return nil, err
	}
	return ret[:],nil
}

type BlockHeader struct {
	BlockRoot 			[]byte
	Signature			[]byte // TODO - checking validity + how many voted?
}

func (header *BlockHeader) Validate(bp *BlockProducer) error {
	sig := &bls.Sign{}
	err := sig.Deserialize(header.Signature)
	if err != nil {
		return err
	}

	if res := sig.VerifyHash(bp.PubKey, header.BlockRoot); !res {
		return fmt.Errorf("signatur did not verify")
	}
	return nil
}