package core

import (
	"crypto/sha256"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/ulule/deepcopier"
)

func CopyState(state *State) *State {
	newBlockRoots := make([]*EpochAndBytes, len(state.BlockRoots))
	for i, r := range state.BlockRoots {
		newBlockRoots[i] = &EpochAndBytes{}
		deepcopier.Copy(r).To(newBlockRoots[i])
	}

	newStateRoots := make([]*EpochAndBytes, len(state.StateRoots))
	for i, r := range state.StateRoots {
		newStateRoots[i] = &EpochAndBytes{}
		deepcopier.Copy(r).To(newStateRoots[i])
	}

	newSeeds := make([]*EpochAndBytes, len(state.Seeds))
	for i, r := range state.Seeds {
		newSeeds[i] = &EpochAndBytes{}
		deepcopier.Copy(r).To(newSeeds[i])
	}

	newBPs := make([]*BlockProducer, len(state.BlockProducers))
	for i, bp := range state.BlockProducers {
		newBPs[i] = &BlockProducer{}
		deepcopier.Copy(bp).To(newBPs[i])
	}

	newPools := make([]*Pool, len(state.Pools))
	for i, p := range state.Pools {
		newPools[i] = &Pool{}
		deepcopier.Copy(p).To(newPools[i])
	}

	newSlashings := make([]uint64, len(state.Slashings))
	for i, s := range state.Slashings {
		deepcopier.Copy(s).To(newSlashings[i])
	}

	return &State{
		GenesisTime:          state.GenesisTime,
		CurrentEpoch:         state.CurrentEpoch,
		BlockRoots:           newBlockRoots,
		StateRoots:           newStateRoots,
		Seeds:                newSeeds,
		BlockProducers:       newBPs,
		Pools:                newPools,
		Slashings:            newSlashings,
	}
}

// will return an 0 length byte array if not found
func GetStateRoot(state *State, epoch uint64) []byte {
	for _, r := range state.StateRoots {
		if r.Epoch == epoch {
			return r.Bytes
		}
	}
	return []byte{}
}

func DecreaseBPBalance(bp *BlockProducer, change uint64) error {
	if bp.Balance < change {
		return fmt.Errorf("BP %d dosen't have enonugh Balance (%d) to decrease (%d)", bp.Id, bp.Balance, change)
	}

	bp.Balance -= change
	return nil
}

func IncreaseBPBalance(bp *BlockProducer, change uint64) error {
	bp.Balance += change
	return nil
}

// will return nil if not found
func GetBlockProducer(state *State, id uint64) *BlockProducer {
	for _, p := range state.BlockProducers {
		if p.GetId() == id {
			return p
		}
	}
	return nil
}

// will return nil if not found
func GetPool(state *State, id uint64) *Pool {
	for _, p := range state.Pools {
		if p.GetId() == id {
			return p
		}
	}
	return nil
}

func PoolCommittee(state *State, poolId uint64, epoch uint64) ([]uint64,error) {
	// TODO - handle integer overflow
	seed, err := GetSeed(state, epoch - 1) // we always use the seed from previous epoch
	if err != nil {
		return []uint64{}, err
	}
	return shuffle(
		state.BlockProducers,
		0,
		epoch,
		shared.SliceToByte32(seed),
		[]byte(fmt.Sprintf("pool %d committee", poolId)),
	)
}

// DKG committee is chosen randomly by shuffling a seed + category (dkg committee)
// The previous epoch's seed is used to choose the DKG committee as the current one (the block's epoch)
func DKGCommittee(state *State, reqId uint64, epoch uint64)([]uint64, error) {
	// TODO - handle integer overflow
	seed, err := GetSeed(state, epoch - 1) // we always use the seed from previous epoch
	if err != nil {
		return []uint64{}, err
	}
	return shuffle(
		state.BlockProducers,
		0,
		epoch,
		shared.SliceToByte32(seed),
		[]byte("dkg committee"),
	)
}

// Block voting committee is chosen randomly by shuffling a seed + category (block voting committee)
// The previous epoch's seed is used to choose the block voting committee as the current one (the block's epoch)
func BlockVotingCommittee(state *State, epoch uint64)([]uint64, error) {
	// TODO - handle integer overflow
	seed, err := GetSeed(state, epoch - 1) // we always use the seed from previous epoch
	if err != nil {
		return []uint64{}, err
	}
	return shuffle(
		state.BlockProducers,
		0,
		epoch,
		shared.SliceToByte32(seed),
		[]byte("block voting committee"),
	)
}

// Block producer is chosen randomly by shuffling a seed + category (block proposer)
// The previous epoch's seed is used to choose the block producer as the current one (the block's epoch)
func GetBlockProposer(state *State, epoch uint64) (uint64, error) {
	seed, err := GetSeed(state, epoch - 1) // we always use the seed from previous epoch
	if err != nil {
		return 0, err
	}

	lst, err := shuffle(
		state.BlockProducers,
		0,
		epoch,
		shared.SliceToByte32(seed),
		[]byte("block proposer"),
	)
	if err != nil {
		return 0, err
	}
	return lst[0], nil
}

// will return error if not found
func GetSeed(state *State, epoch uint64) ([]byte, error) {
	for _, d := range state.Seeds {
		if d.Epoch == epoch {
			return d.Bytes, nil
		}
	}
	return []byte{}, fmt.Errorf("seed for epoch %d not found", epoch)
}

// TODO - find out if secure
func shuffle(allBPs []*BlockProducer, committeeId uint64, epoch uint64, seed [32]byte, nonce []byte) ([]uint64, error) {
	// get Active BPs
	var activeBps []uint64
	for _, bp := range allBPs {
		if bp.Active || bp.GetExitEpoch() > epoch {
			activeBps = append(activeBps, bp.GetId())
		}
	}

	// nonce is used as different categories for the seed
	seedToUse := seed
	if nonce != nil {
		h := sha256.New() // TODO - secure enough?
		_, err := h.Write(append(seed[:], nonce...))
		if err != nil {
			return []uint64{}, err
		}
		seedToUse = shared.SliceToByte32(h.Sum(nil))
	}

	// shuffle
	shuffled,err := shared.ShuffleList(activeBps, seedToUse, 60)
	if err != nil {
		return nil, err
	}

	//
	ret := make([]uint64, TestConfig().PoolExecutorsNumber)
	for i := uint64(0) ; i < TestConfig().PoolExecutorsNumber ; i++ {
		ret[i] = shuffled[committeeId* TestConfig().PoolExecutorsNumber + i]
	}

	return ret, nil
}