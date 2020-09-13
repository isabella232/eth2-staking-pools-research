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

	newETH1 := make([]*EpochAndNumber, len(state.ETH1Blocks))
	for i, s := range state.ETH1Blocks {
		newETH1[i] = &EpochAndNumber{}
		deepcopier.Copy(s).To(newETH1[i])
	}

	newETH2 := make([]*EpochAndNumber, len(state.ETH2Epochs))
	for i, s := range state.ETH2Epochs {
		newETH2[i] = &EpochAndNumber{}
		deepcopier.Copy(s).To(newETH2[i])
	}

	return &State{
		GenesisTime:    state.GenesisTime,
		CurrentEpoch:   state.CurrentEpoch,
		BlockRoots:     newBlockRoots,
		StateRoots:     newStateRoots,
		Seeds:          newSeeds,
		BlockProducers: newBPs,
		Pools:          newPools,
		Slashings:      newSlashings,
		ETH1Blocks:     newETH1,
		ETH2Epochs:     newETH2,
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

func GetActiveBlockProducers(state *State, epoch uint64) []uint64 {
	var activeBps []uint64
	for _, bp := range state.BlockProducers {
		if bp.Active || bp.GetExitEpoch() > epoch {
			activeBps = append(activeBps, bp.GetId())
		}
	}
	return activeBps
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

// Pool committee is a randomly selected committee of BPs that are chosen to receive a pool's
// keys via key rotation.
// For new pools, the first committee is responsible for doing DKG
//
// Pool committee is chosen randomly by shuffling a seed + category (pool %d committee)
// The previous epoch's seed is used to choose the DKG committee as the current one (the block's epoch)
func PoolCommittee(state *State, poolId uint64, epoch uint64) ([]uint64,error) {
	// TODO - handle integer overflow
	seed, err := GetSeed(state, epoch - 1) // we always use the seed from previous epoch
	if err != nil {
		return []uint64{}, err
	}
	return shuffleActiveBPs(
		GetActiveBlockProducers(state, epoch),
		shared.SliceToByte32(seed),
		[]byte(fmt.Sprintf("pool %d committee", poolId)),
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
	return shuffleActiveBPs(
		GetActiveBlockProducers(state, epoch),
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

	lst, err := shuffleActiveBPs(
		GetActiveBlockProducers(state, epoch),
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

// Shuffle takes in a list of block producers Ids, a seed and a nonce to create a unique shuffle for that
// combination by hashing seed + nonce.
// Changing the nonce for different purposes can be used as "categories" from the same seed
// TODO - find out if secure
func shuffleActiveBPs(bps []uint64, seed [32]byte, nonce []byte) ([]uint64, error) {
	// nonce is used to randomly select multiple types of committees from the same seed
	seedToUse := seed
	if nonce != nil {
		h := sha256.New() // TODO - secure enough?
		_, err := h.Write(append(seed[:], nonce...))
		if err != nil {
			return []uint64{}, err
		}
		seedToUse = shared.SliceToByte32(h.Sum(nil))
	}

	// shuffleActiveBPs
	shuffled,err := shared.ShuffleList(bps, seedToUse, 60)
	if err != nil {
		return nil, err
	}

	//
	ret := make([]uint64, TestConfig().PoolExecutorsNumber)
	for i := uint64(0) ; i < TestConfig().PoolExecutorsNumber ; i++ {
		ret[i] = shuffled[TestConfig().PoolExecutorsNumber + i]
	}

	return ret, nil
}