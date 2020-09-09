package block

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
)

// A request struct for creating new pool credentials
// will trigger random selection of 128 executors to DKG new pool credentials and wait for deposit
//
// How it works?
// - A user sends 32 eth and create pool request
// - The first BP that sees it, will post a CreatePoolRequest with Status 0 and will nominate the next BP as the leader for the DKG
//   (the 128 DKG participants are deterministically selected as well)
// - If during the next Epoch the DKG is successful, the BP (which is also the DKG leader) posts a CreatePoolRequest with the same ID,
//   Status 1 and the created pub key
// - If the DKG is un-successful, the BP will post a CreatePoolRequest with the same ID, Status 3 and will nominate the next BP as leader
//
// A successful DKG will reward the leader and DKG participants
// A non-successful DKG will penalized the DKG participants
type CreatePoolRequest struct {
	Id                  uint64 // primary key
	Status              uint64 // 0 for not completed, 1 for completed, 2 for un-successful
	StartEpoch          uint64
	EndEpoch            uint64
	LeaderBlockProducer uint64   // should be the next block producer
	CreatedPubKey       []byte   // populated after DKG is successful
	Participation       [16]byte // 128 bit of the executors (by order) which executed this duty
}

func NewCreatePoolRequest(
	id					uint64,
	status 				uint64,
	startEpoch			uint64,
	endEpoch			uint64,
	leaderBlockProducer	uint64,
	createdPubKey		[]byte,
	participation		[16]byte,
	) *CreatePoolRequest {
	return &CreatePoolRequest{
		Id:                  id,
		Status:              status,
		StartEpoch:          startEpoch,
		EndEpoch:            endEpoch,
		LeaderBlockProducer: leaderBlockProducer,
		CreatedPubKey:       createdPubKey,
		Participation:       participation,
	}
}

func (req *CreatePoolRequest) GetId() uint64 {
	return req.Id
}

func (req *CreatePoolRequest) GetStatus() uint64 {
	return req.Status
}

func (req *CreatePoolRequest) GetStartEpoch() uint64 {
	return req.StartEpoch
}

func (req *CreatePoolRequest) GetEndEpoch() uint64 {
	return req.EndEpoch
}

func (req *CreatePoolRequest) GetLeaderBP() uint64 {
	return req.LeaderBlockProducer
}

func (req *CreatePoolRequest) GetCreatePubKey() []byte {
	return req.CreatedPubKey
}

func (req *CreatePoolRequest) GetParticipation() [16]byte  {
	return req.Participation
}

func (req *CreatePoolRequest) Validate(state core.IState, currentBP core.IBlockProducer) error {
	if req.LeaderBlockProducer != currentBP.GetId() {
		return fmt.Errorf("pool leader should be the current block producer")
	}

	// TODO - req Id is primary (non duplicate and incremental)

	// TODO - check that network has enough capitalization

	// TODO - check leader is not part of DKG Committee
	return nil
}
