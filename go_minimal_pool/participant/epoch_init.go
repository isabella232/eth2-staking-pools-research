package participant

import (
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/crypto"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net/pb"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/state"
	"github.com/google/uuid"
	"log"
)

// start happens at 1/4 of the epoch
// https://github.com/bloxapp/eth2-staking-pools-research/blob/master/epoch_processing.md
func (p *Participant) epochInit(epoch *state.Epoch) {
	p.epochProcessingLock.Lock()
	defer p.epochProcessingLock.Unlock()

	log.Printf("P %d, epoch %d init", p.Id, epoch.Number)

	// find share distro target
	nextEpoch := p.Node.State.GetEpoch(epoch.Number + 1)
	nextEpochPools,err := nextEpoch.PoolsParticipantIds()
	if err != nil {
		log.Fatalf("P %d err fetching next epoch's pools: %s", p.Id, err.Error())
	}
	currentPool,err := epoch.ParticipantPoolAssignment(p.Id)
	if err != nil {
		log.Fatalf("P %d err fetching current epoch's pool: %s", p.Id, err.Error())
	}
	sharePoolTarget := nextEpochPools[currentPool]


	// generate re-distro shares
	config := net.NewTestNetworkConfig()
	distro,err := crypto.NewRedistribuition(config.PoolThreshold -1, epoch.ParticipantShare)
	if err != nil {
		log.Fatalf("P %d err instantiating NewRedistribuition: %s", p.Id, err.Error())
	}
	shares,err := distro.GenerateShares(sharePoolTarget)
	if err != nil {
		log.Fatalf("P %d err generating re-distro shares: %s", p.Id, err.Error())
	}

	// broadcast
	for k,v := range shares {
		share := &pb.ShareDistribution{
			Id:              uuid.New().String(),
			FromParticipant: &pb.Participant{Id: p.Id},
			ToParticipant:   &pb.Participant{Id: k},
			Share:           v.Serialize(),
			Commitments:     nil,
			PoolId:          uint32(currentPool),
			Epoch:           epoch.Number,
		}

		err := p.Node.Net.BroadcastShare(share)
		if err != nil {
			log.Printf("broadcasting error: %s", err.Error())
		}
	}
}