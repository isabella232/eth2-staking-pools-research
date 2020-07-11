package participant

import (
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/crypto"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net/pb"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/state"
	"github.com/google/uuid"
	"log"
)

// start happens at 1/2 of the epoch
// https://github.com/bloxapp/eth2-staking-pools-research/blob/master/epoch_processing.md
func (p *Participant) epochMid(epoch *state.Epoch) {
	log.Printf("P %d, epoch %d mid with %d shares", p.Id, epoch.Number, len(p.Node.SharesPerEpoch[epoch.Number]))

	currentPool,err := epoch.ParticipantPoolAssignment(p.Id)
	if err != nil {
		log.Fatalf("P %d err fetching current epoch's pool: %s", p.Id, err.Error())
	}

	config := net.NewTestNetworkConfig()
	sigInG2 := crypto.Sign(epoch.ParticipantShare, config.EpochTestMessage)
	sig := &pb.SignatureDistribution{
		Id:              uuid.New().String(),
		FromParticipant: &pb.Participant{Id: p.Id},
		Sig:           	 sigInG2.Serialize(),
		PoolId:          uint32(currentPool),
		Epoch:           epoch.Number,
	}
	err = p.Node.Net.BroadcastSignature(sig)
	if err != nil {
		log.Printf("broadcasting error: %s", err.Error())
	}
}
