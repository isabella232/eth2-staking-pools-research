package main

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/crypto"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/participant"
	pool_chain "github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/state"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
	"log"
)

var participants []*participant.Participant

func main() {
	crypto.InitBLS()
	log.SetFlags(log.Lmicroseconds)

	config := net.NewTestNetworkConfig()
	participants = make([]*participant.Participant, config.TotalNumberOfParticipants())

	// simulate epoch 0 and get initial pool assignments
	seed,err := crypto.MixSeed(config.GenesisSeed, 0)
	if err != nil {
		log.Fatalf(err.Error())
	}
	epoch := state.NewEpochInstance(0, seed)
	poolData,err := epoch.PoolsParticipantIds()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// initial DKG for pools
	participants = runDKGForPools(poolData, config.PoolThreshold - 1)

	// connect pools to each other
	for i, p1 := range participants {
		p1.Node.Net.AddPeer(p1.Node.Net.OwnPeer()) // add to self to receive shares
		for j := i+1 ; j < len(participants) ; j ++  {
			p2 := participants[j]
			net.BiDirectionalConnection(p1.Node.Net,p2.Node.Net)
		}
	}
	// start epoch processing
	for _, p := range participants {
		p.StartEpochProcessing()
	}

	for {
		select {
		case sig := <- participants[0].KillC():
			if sig == true {
				fmt.Printf("killed")
				return
			}
		}
	}
}

// will run a DKG for every pool inputed, creates the participant and it's node, generated pool shared secret
// and sets pool data as well
func runDKGForPools(poolData map[shared.PoolId][]shared.ParticipantId, threshold shared.PoolSize) []*participant.Participant {
	ret := make([]*participant.Participant,0)
	pools := make([]*state.Pool, len(poolData))
	i := 0
	for poolId, poolParticipants := range poolData {
		sks, pk, err := runDKGForParticipants(threshold - 1, poolParticipants)
		if err != nil {
			log.Fatalf(err.Error())
		}

		log.Printf("pool %d:", poolId)
		for k, v := range sks {
			log.Printf("		p %d, sk: %s", k, v.GetString(10))

			// create the participant
			p := participant.NewParticipant(k)
			n := pool_chain.NewTestChainNode()
			p.SetNode(n)
			// set secret
			e := n.State.GetEpoch(0)
			e.ParticipantShare = v
			n.State.SaveEpoch(e)

			ret = append(ret, p)
		}

		// create pool data
		pools[i] = state.NewPool(poolId, threshold, pk)
		i++
	}

	// for each participant save pool data
	for _, p := range ret {
		for _, pool := range pools {
			p.Node.State.SavePool(pool)
		}
	}

	return ret
}

func runDKGForParticipants(degree shared.PoolSize, indexes []shared.ParticipantId) (map[shared.ParticipantId]*bls.Fr, *bls.PublicKey, error) {
	dkg,err := crypto.NewDKG(degree, indexes)
	if err != nil {
		return nil, nil, err
	}

	sks, err := dkg.GroupSecrets(indexes)
	if err != nil {
		return nil, nil, err
	}
	pk,err := dkg.GroupPK(sks)
	if err != nil {
		return nil, nil, err
	}

	return sks, pk, nil
}
