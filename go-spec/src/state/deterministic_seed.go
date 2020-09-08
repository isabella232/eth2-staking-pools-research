package state

import (
	"crypto/sha1"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
)

// TODO - should be randomly choosen depending on epoch
func (state *State) PoolExecutors(poolId uint64, epoch uint64) ([]uint64, error) {
	//if !state.IsActivePool(poolId) {
	//	return nil, fmt.Errorf("pool not active")
	//}

	// get active BPs
	var activeBps []uint64
	for _, bp := range state.blockProducers {
		if bp.IsActive() {
			activeBps = append(activeBps, bp.GetId())
		}
	}

	shuffled,err := shared.ShuffleList(activeBps, state.GetSeed(state.GetCurrentEpoch()), 10)
	if err != nil {
		return nil, err
	}

	ret := make([]uint64, core.TestConfig().PoolExecutorsNumber)
	for i := uint64(0) ; i < core.TestConfig().PoolExecutorsNumber ; i++ {
		ret[i] = shuffled[poolId * core.TestConfig().PoolExecutorsNumber + i]
	}

	return ret, nil
}

// Deterministic method to randomly pick DKG participants based on state seed and req number
// TODO - find better way
func (state *State) DKGCommittee(reqId uint64, epoch uint64)([]uint64, error) {
	ret := make([]uint64, core.TestConfig().PoolExecutorsNumber)

	// get active BPs
	var activeBps []uint64
	for _, bp := range state.blockProducers {
		if bp.IsActive() {
			activeBps = append(activeBps, bp.GetId())
		}
	}

	h := sha1.New()
	seed := state.GetSeed(state.GetCurrentEpoch())
	for i := 0 ; i < int(reqId) % 100 ; i++ {
		_, err := h.Write(seed[:])
		if err != nil {
			return ret, err
		}

		seed = shared.SliceToByte32(h.Sum(nil))
	}

	shuffled,err := shared.ShuffleList(activeBps, seed, 20)
	if err != nil {
		return nil, err
	}


	for i := uint64(0) ; i < core.TestConfig().DKGParticipantsNumber ; i++ {
		ret[i] = shuffled[i]
	}

	return ret, nil
}

func (state *State) BlockVotingCommittee(epoch uint64)([]uint64, error) {
	return []uint64{},nil
}

func (state *State) GetBlockProposer(epoch uint64) (uint64, error) {
	return 0, nil
}

