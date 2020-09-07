package src

import (
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-ssz"
)

// A request struct for creating new pool credentials
// will trigger random selection of 128 executors to DKG new pool credentials and wait for deposit
//
// How it works?
// - A user sends 32 eth and create pool request
// - The first BP that sees it, will post a CreatePoolRequest with status 0 and will nominate the next BP as the leader for the DKG
//   (the 128 DKG participants are deterministically selected as well)
// - If during the next epoch the DKG is successful, the BP (which is also the DKG leader) posts a CreatePoolRequest with the same ID,
//   status 1 and the created pub key
// - If the DKG is un-successful, the BP will post a CreatePoolRequest with the same ID, status 3 and will nominate the next BP as leader
//
// A successful DKG will reward the leader and DKG participants
// A non-successful DKG will penalized the DKG participants
type CreatePoolRequest struct {
	Id					uint64 // primary key
	Status 				uint64 // 0 for not completed, 1 for completed, 2 for un-successful
	StartEpoch			uint64
	EndEpoch			uint64
	LeaderBlockProducer	uint64 // should be the next block producer
	CreatedPubKey		[]byte // populated after DKG is successful
	Participation		[16]byte // 128 bit of the executors (by order) which executed this duty
}
func (req *CreatePoolRequest) Validate(state *State, currentBP *BlockProducer) error {
	if req.LeaderBlockProducer != currentBP.Id {
		return fmt.Errorf("pool leader should be the current block producer")
	}

	// TODO - req id is primary (non duplicate and incremental)

	// TODO - check that network has enough capitalization
	return nil
}


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
	NewPoolReq				[]*CreatePoolRequest
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

	headBlockBody := helperFunc.GetBlockBody(state.HeadBlockHeader.BlockRoot)
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

func (header *BlockHeader) Copy() *BlockHeader {
	return &BlockHeader{
		BlockRoot: header.BlockRoot,
		Signature: header.Signature,
	}
}

func NewBlockHeader(sk *bls.SecretKey, body *BlockBody) (*BlockHeader,error) {
	root,err := body.Root()
	if err != nil {
		return nil, err
	}

	sig := sk.SignByte(root)

	return &BlockHeader{
		BlockRoot: root,
		Signature: sig.Serialize(),
	}, nil
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