package state

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/crypto"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
)

func shufflePools(input []shared.ParticipantId, seed [32]byte, roundCount uint8, numberOfPools shared.PoolId, poolSize shared.PoolSize) (map[shared.PoolId][]shared.ParticipantId, error) {
	shuffled, err := crypto.ShuffleList(input, seed, roundCount)
	if err != nil {
		return nil, err
	}

	ret := make(map[shared.PoolId][]shared.ParticipantId)
	for p_id := shared.PoolId(1) ; p_id <= numberOfPools ; p_id ++ {
		start := int(p_id - 1) * int(poolSize)
		end := start + int(poolSize)
		ret[p_id] = shuffled[start: end]
	}

	return ret, nil
}

type Epoch struct {
	Number shared.EpochNumber
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

func (epoch *Epoch) ParticipantPoolAssignment(id shared.ParticipantId) (shared.PoolId,error) {
	// TODO make more efficient
	pools,err := epoch.PoolsParticipantIds()
	if err != nil {
		return 0,err
	}

	for poolId, pool := range pools {
		for _, _pId := range pool {
			if _pId == id {
				return poolId, nil
			}
		}
	}

	return 0,fmt.Errorf("can't find %d", id)
}

func (epoch *Epoch) PoolsParticipantIds() (map[shared.PoolId][]shared.ParticipantId,error) {
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