package state

import (
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/crypto"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net"
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
}

func NewEpochInstance(number uint32, seed [32]byte) *Epoch {
	return &Epoch{
		Number:number,
		epochSeed: seed,
	}
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
