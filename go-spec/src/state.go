package src

import (
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-ssz"
)

var helperFunc NonSpecFunctions

type BlockProducer struct {
	Id				uint64
	PubKey			*bls.PublicKey
	Balance			uint64 // balance on the pool chain (rewards earned)
	Stake			uint64 // stake
	Slashed			bool
	Active 			bool
}

func (bp *BlockProducer) Copy() (*BlockProducer, error) {
	pk := &bls.PublicKey{}
	err := pk.Deserialize(bp.PubKey.Serialize())
	if err != nil {
		return nil, err
	}

	return &BlockProducer{
		Id:      bp.Id,
		PubKey:  pk,
		Balance: bp.Balance,
		Stake:   bp.Stake,
		Slashed: bp.Slashed,
		Active:  bp.Active,
	}, nil
}

type Pool struct {
	Id					uint64 // id
	PubKey				*bls.PublicKey // eth2 validation pubkey
	SortedExecutors		[]uint64 // ids of the block producers which are executors on this pool
}
func (pool *Pool) Copy() (*Pool, error) {
	pk := &bls.PublicKey{}
	err := pk.Deserialize(pool.PubKey.Serialize())
	if err != nil {
		return nil, err
	}

	return &Pool{
		Id:              pool.Id,
		PubKey:          pk,
		SortedExecutors: pool.SortedExecutors,
	}, nil
}

type State struct {
	Pools			[]*Pool
	BlockRoots		[]byte
	HeadBlockHeader	*BlockHeader
	BlockProducers  []*BlockProducer
	Seed			[32]byte
}

func (state *State) Copy() (*State, error) {
	copiedPools := make([]*Pool, len(state.Pools))
	for i, p := range state.Pools {
		newP, err := p.Copy()
		if err != nil {
			return nil, err
		}
		copiedPools[i] = newP
	}

	copiedBps := make([]*BlockProducer, len(state.BlockProducers))
	for i, bp := range state.BlockProducers {
		newBP, err := bp.Copy()
		if err != nil {
			return nil, err
		}
		copiedBps[i] = newBP
	}

	return &State{
		Pools:           copiedPools,
		BlockRoots:      state.BlockRoots,
		HeadBlockHeader: state.HeadBlockHeader.Copy(),
		BlockProducers:  copiedBps,
		Seed:            state.Seed,
	}, nil
}

func (state *State) Root() ([32]byte,error) {
	return ssz.HashTreeRoot(state)
}

func (state *State) IsActivePool(id uint64) bool {
	return true // TODO
}

func (state *State)  GetPool(id uint64) *Pool {
	for _, p := range state.Pools {
		if p.Id == id {
			return p
		}
	}
	return nil
}

func (state *State)  GetBlockProducer(id uint64) (*BlockProducer, error) {
	for _, bp := range state.BlockProducers {
		if bp.Id == id {
			return bp, nil
		}
	}
	return nil, fmt.Errorf("block producer not found")
}

// TODO - should be randomly choosen depending on epoch
func (state *State) PoolExecutors(poolId uint64, epoch uint64) ([]uint64, error) {
	if !state.IsActivePool(poolId) {
		return nil, fmt.Errorf("pool not active")
	}

	// get active BPs
	var activeBps []uint64
	for _, bp := range state.BlockProducers {
		if bp.Active {
			activeBps = append(activeBps, bp.Id)
		}
	}

	shuffled,err := ShuffleList(activeBps, state.Seed, 10)
	if err != nil {
		return nil, err
	}

	ret := make([]uint64, TestConfig().PoolExecutorsNumber)
	for i := uint64(0) ; i < TestConfig().PoolExecutorsNumber ; i++ {
		ret[i] = shuffled[poolId * TestConfig().PoolExecutorsNumber + i]
	}

	return ret, nil
}

func (state *State) GetBlockProposer(epoch uint64) uint64 {
	return 0
}

func (state *State) IncreaseBlockProducerBalance(id uint64, change uint64) (newBalance uint64, error error) {
	bp,err := state.GetBlockProducer(id)
	if err != nil {
		return 0, err
	}

	bp.Balance += change
	return bp.Balance, nil
}

func (state *State) DecreaseBlockProducerBalance(id uint64, change uint64) (newBalance uint64, error error) {
	bp,err := state.GetBlockProducer(id)
	if err != nil {
		return 0, err
	}

	if bp.Balance < change {
		return 0, fmt.Errorf("BP %d dosen't have enonugh balance (%d) to decrease (%d)", bp.Id, bp.Balance, change)
	}

	bp.Balance -= change
	return bp.Balance, nil
}

func (state *State) ValidateBlock(header *BlockHeader, body *BlockBody) error {
	bp, err := state.GetBlockProducer(body.Proposer)
	if err != nil {
		return err
	}

	err = header.Validate(bp)
	if err != nil {
		return err
	}

	err = body.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (state *State) AddNewPool(pool *Pool) error {
	if found := state.GetPool(pool.Id); found != nil {
		return fmt.Errorf("pool already exists")
	}

	state.Pools = append(state.Pools, pool)
	return nil
}

// Applies every pool performance to its relevant executors, decreasing and increasing balances.
func (state *State) ApplyPoolExecutions(summaries []*PoolExecutionSummary) error {
	for _, summary := range summaries {
		if !state.IsActivePool(summary.PoolId) {
			return fmt.Errorf("pool %d not active", summary.PoolId)
		}

		if err := summary.ApplyOnState(state); err != nil {
			return err
		}
	}
	return nil
}

//func (state *State) ProcessNewPoolRequests(requests []*CreatePoolRequest, currentBP *BlockProducer) error {
//	//for _, req := range requests {
//	//
//	//}
//	return nil
//}

// called when a new block was proposed
func (state *State) ProcessNewBlock(newBlockHeader *BlockHeader, newBlockBody *BlockBody) (newState *State, error error) {
	if state.GetBlockProposer(newBlockBody.Number) != newBlockBody.Proposer {
		return nil, fmt.Errorf("block proposer is worng")
	}


	// copy the state to apply state transition on
	stateCopy, err := state.Copy()
	if err != nil {
		return nil, err
	}

	err = stateCopy.ApplyPoolExecutions(newBlockBody.PoolsExecutionSummary)
	if err != nil {
		return nil, err
	}

	stateCopy.Seed, err = MixSeed(stateCopy.Seed, SliceToByte32(newBlockHeader.Signature[:32])) // TODO - use something else than the sig
	if err != nil {
		return nil, err
	}

	err = helperFunc.SaveBlockBody(newBlockBody)
	if err != nil {
		return nil, err
	}

	return stateCopy, nil
}