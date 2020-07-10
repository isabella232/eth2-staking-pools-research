package state

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/crypto"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net"
	"github.com/herumi/bls-eth-go-binary/bls"
)

func shufflePools(input []uint32, seed [32]byte, roundCount uint8, numberOfPools uint8, poolSize uint8) (map[uint8][]uint32, error) {
	shuffled, err := crypto.ShuffleList(input, seed, roundCount)
	if err != nil {
		return nil, err
	}

	ret := make(map[uint8][]uint32)
	for p_id := uint8(1) ; p_id <= numberOfPools ; p_id ++ {
		start := (p_id - 1) * poolSize
		end := start + poolSize
		ret[p_id] = shuffled[start: end]
	}

	return ret, nil
}

type Epoch struct {
	Number uint32
	epochSeed [32]byte

	// every participant will use this var to store his epoch's secret.
	ParticipantShare *bls.Fr
	// used to store the epoch's reconstructed signature (that will get broadcasted to eth2)
	ReconstructedSignature *bls.G2
	//
	EpochSigVerified bool
}

func NewEpochInstance(number uint32, seed [32]byte) *Epoch {
	return &Epoch{
		Number:number,
		epochSeed: seed,
		EpochSigVerified: false,
	}
}

func (epoch *Epoch) ParticipantPoolAssignment(participantId uint32) (uint8,error) {
	// TODO make more efficient
	pools,err := epoch.PoolsParticipantIds()
	if err != nil {
		return 0,err
	}

	for poolId, pool := range pools {
		for _, _pId := range pool {
			if _pId == participantId {
				return poolId, nil
			}
		}
	}

	return 0,fmt.Errorf("can't find %d", participantId)
}

func (epoch *Epoch) PoolsParticipantIds() (map[uint8][]uint32,error) {
	config := net.NewTestNetworkConfig()
	return shufflePools(
		config.ParticipantIndexesList(),
		epoch.epochSeed,
		config.SeedShuffleRoudnCount,
		config.NumberOfPools,
		config.PoolSize,
		)
}

func (epoch *Epoch)StatusString() string {
	return fmt.Sprintf("Epoch number: %d, Sig Verified: %t",epoch.Number,epoch.EpochSigVerified)
}