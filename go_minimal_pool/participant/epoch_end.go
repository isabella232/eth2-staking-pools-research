package participant

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/crypto"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/state"
	"github.com/herumi/bls-eth-go-binary/bls"
	"log"
)

// start happens at 2/3 of the epoch
// https://github.com/bloxapp/eth2-staking-pools-research/blob/master/epoch_processing.md
func (p *Participant) epochEnd(epoch *state.Epoch) {
	p.epochProcessingLock.Lock()
	defer p.epochProcessingLock.Unlock()

	log.Printf("P %d, epoch %d end with %d sigs", p.Id,epoch.Number, len(p.Node.SigsPerEpoch[epoch.Number]))

	err := p.reconstructEpochSignature(epoch)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	err = p.verifyEpochSig(epoch)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	err = p.reconstructGroupSecretForNextEpoch(epoch)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}


	currentPool,_ := epoch.ParticipantPoolAssignment(p.Id)
	log.Printf("P %d, pool: %d, epoch status: %s",p.Id,currentPool, epoch.StatusString())
}


func (p *Participant) verifyEpochSig(epoch *state.Epoch) error {
	config := net.NewTestNetworkConfig()
	currentPool,err := epoch.ParticipantPoolAssignment(p.Id)
	if err != nil {
		return fmt.Errorf("P %d err fetching current epoch's pool: %s", p.Id, err.Error())
	}

	sig := bls.CastToSign(epoch.ReconstructedSignature)
	pk := p.Node.State.Pools[currentPool].Pk
	epoch.EpochSigVerified = sig.VerifyByte(pk, config.EpochTestMessage)
	p.Node.State.SaveEpoch(epoch)
	return nil
}

func (p *Participant) reconstructEpochSignature(epoch *state.Epoch) error {
	currentPool,err := epoch.ParticipantPoolAssignment(p.Id)
	if err != nil {
		return fmt.Errorf("P %d err fetching current epoch's pool: %s", p.Id, err.Error())
	}

	// filter out relevant sigs
	points := make([][]interface{},0)
	for _,v := range p.Node.SigsPerEpoch[epoch.Number] {
		if v.PoolId == currentPool {
			sig := &bls.G2{}
			err := sig.Deserialize(v.Sig)
			if err != nil {
				continue
			}

			id := &bls.Fr{}
			id.SetInt64(int64(v.FromParticipant.Id))

			points = append(points, []interface{}{*id, sig})
		}
	}

	// reconstruct
	l := crypto.NewG2LagrangeInterpolation(points)
	sig,err := l.Interpolate()
	if err != nil {
		return fmt.Errorf("could not reconstruct group signature for epoch %d: %s", epoch.Number, err.Error())
	}

	epoch.ReconstructedSignature = sig
	p.Node.State.SaveEpoch(epoch)
	return nil
}

func (p *Participant) reconstructGroupSecretForNextEpoch(epoch *state.Epoch) error {
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
		return fmt.Errorf("could not reconstruct group secret for next epoch: %s", err.Error())
	}

	// save for next epoch
	nextEpoch := p.Node.State.GetEpoch(epoch.Number + 1)
	nextEpoch.ParticipantShare = groupSk
	err = p.Node.State.SaveEpoch(nextEpoch)
	if err != nil {
		return fmt.Errorf("could not save group secret for next epoch: %s", err.Error())
	}
	return nil
}