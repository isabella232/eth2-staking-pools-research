package participant

import (
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/crypto"
	pool_chain "github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net/pb"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/state"
	"github.com/google/uuid"
	"github.com/herumi/bls-eth-go-binary/bls"
	"log"
	"sync"
	"time"
)

type Participant struct {
	Id   uint32
	Node *pool_chain.PoolChainNode

	epochProcessingLock sync.Mutex
}

func NewParticipant(id uint32) *Participant {
	return &Participant{
		Id:   id,
	}
}

func (p *Participant) SetNode(node *pool_chain.PoolChainNode) {
	p.Node = node
	p.Node.FilterId = p.Id
}

func (p *Participant) StartEpochProcessing() {
	go func() {
		for {
			select {
			case epoch := <- p.Node.EpochC():
				if epoch == -1 {
					return
				}

				e := p.Node.State.GetEpoch(uint32(epoch))
				go p.timeEpoch(e)
			}
		}
	}()

	p.Node.StartEpochProcessing()

	log.Printf("Participant %d started", p.Id)
}

func (p *Participant) KillC() <- chan bool {
	return p.Node.Killed
}

func (p *Participant) timeEpoch(epoch *state.Epoch) {
	// start happens at 1/4 of the epoch
	go func() {
		d := time.Duration(p.Node.Config.EpochSpanSec / 4)
		<- time.After(d)

		p.epochInit(epoch)
	}()

	// mid happens at 1/2 of the epoch
	go func() {
		d := time.Duration(p.Node.Config.EpochSpanSec / 2)
		<- time.After(d)

		p.epochMid(epoch)
	}()

	// end happens at 3/4 of the epoch
	go func() {
		d := time.Duration((p.Node.Config.EpochSpanSec / 4) * 3)
		<- time.After(d)

		p.epochEnd(epoch)
	}()
}

// start happens at 1/4 of the epoch
// https://github.com/bloxapp/eth2-staking-pools-research/blob/master/epoch_processing.md
func (p *Participant) epochInit(epoch *state.Epoch) {
	log.Printf("P %d, epoch %d init", p.Id, epoch.Number)

	p.epochProcessingLock.Lock()
	defer p.epochProcessingLock.Unlock()

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

// start happens at 2/3 of the epoch
// https://github.com/bloxapp/eth2-staking-pools-research/blob/master/epoch_processing.md
func (p *Participant) epochEnd(epoch *state.Epoch) {
	log.Printf("P %d, epoch %d end with %d sigs", p.Id,epoch.Number, len(p.Node.SigsPerEpoch[epoch.Number]))

	p.reconstructGroupSecretForNextEpoch(epoch)
}

func (p *Participant) reconstructGroupSecretForNextEpoch(epoch *state.Epoch) {
	shares := p.Node.SharesPerEpoch[epoch.Number]
	points := make([][]bls.Fr,0)
	for _,v := range shares {
		if v.ToParticipant.Id != p.Id {
			continue
		}

		from := &bls.Fr{}
		from.SetInt64(int64(v.FromParticipant.Id))

		point := &bls.Fr{}
		point.Deserialize(v.Share)

		points = append(points, []bls.Fr{*from, *point})
	}

	// reconstruct the group secret from the shares
	l := crypto.NewLagrangeInterpolation(points)
	groupSk, err := l.Interpolate()
	if err != nil {
		log.Printf("could not reconstruct group secret for next epoch: %s", err.Error())
	}

	// save for next epoch
	nextEpoch := p.Node.State.GetEpoch(epoch.Number + 1)
	nextEpoch.ParticipantShare = groupSk
	err = p.Node.State.SaveEpoch(nextEpoch)
	if err != nil {
		log.Printf("could not save group secret for next epoch: %s", err.Error())
	}
}