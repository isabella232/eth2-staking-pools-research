package src

import (
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
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
	Proposer 			uint64
	PoolsExecutionSummary []*PoolExecutionSummary
	//NewPoolReq			[]*CreatePoolRequest
	//WithdrawReq			[]*WithdrawRequest
	//LiquidationReq		[]*LiquidatePoolRequest
	//Slashing			[]*Slashing
	StateRoot			[]byte
	ParentBlockRoot		[]byte
}

func (header *BlockBody) Validate() error {
	return nil
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