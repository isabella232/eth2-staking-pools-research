package state

import (
	"crypto/sha1"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src"
)

// TODO - should be randomly choosen depending on epoch
func (state *State) PoolExecutors(poolId uint64, epoch uint64) ([]uint64, error) {
	//if !state.IsActivePool(poolId) {
	//	return nil, fmt.Errorf("pool not active")
	//}

	// get active BPs
	var activeBps []uint64
	for _, bp := range state.blockProducers {
		if bp.Active {
			activeBps = append(activeBps, bp.Id)
		}
	}

	shuffled,err := src.ShuffleList(activeBps, state.seed, 10)
	if err != nil {
		return nil, err
	}

	ret := make([]uint64, src.TestConfig().PoolExecutorsNumber)
	for i := uint64(0) ; i < src.TestConfig().PoolExecutorsNumber ; i++ {
		ret[i] = shuffled[poolId * src.TestConfig().PoolExecutorsNumber + i]
	}

	return ret, nil
}

// Deterministic method to randomly pick DKG participants based on state seed and req number
// TODO - find better way
func (state *State) DKGCommittee(reqId uint64, epoch uint64)([]uint64, error) {
	ret := make([]uint64, src.TestConfig().PoolExecutorsNumber)

	// get active BPs
	var activeBps []uint64
	for _, bp := range state.blockProducers {
		if bp.Active {
			activeBps = append(activeBps, bp.Id)
		}
	}

	h := sha1.New()
	seed := state.seed[:]
	for i := 0 ; i < int(reqId) % 100 ; i++ {
		_, err := h.Write(seed)
		if err != nil {
			return ret, err
		}

		seed = h.Sum(nil)
	}

	shuffled,err := src.ShuffleList(activeBps, src.SliceToByte32(seed), 10)
	if err != nil {
		return nil, err
	}


	for i := uint64(0) ; i < src.TestConfig().DKGParticipantsNumber ; i++ {
		ret[i] = shuffled[i]
	}

	return ret, nil
}

func (state *State) BlockVotingCommittee(epoch uint64)([]uint64, error) {

}

func (state *State) GetBlockProposer(epoch uint64) (uint64, error) {
	return 0
}

