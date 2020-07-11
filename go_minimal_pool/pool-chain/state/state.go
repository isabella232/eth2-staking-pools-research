package state

import (
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/crypto"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/shared"
)

type DB interface {
	// will return nil,nil if epoch not found
	GetEpoch(number shared.EpochNumber) (*Epoch,error)
	SaveEpoch(epoch *Epoch) error
}

type State struct {
	db           DB
	Pools        map[shared.PoolId]*Pool
	seed         [32]byte
}

func NewInMemoryState(seed [32]byte) *State {
	return & State{
		db:           NewInMemoryDb(),
		Pools:        make(map[shared.PoolId]*Pool),
		seed:         seed,
	}
}

func (s *State) SaveEpoch(epoch *Epoch) error {
	return s.db.SaveEpoch(epoch)
}

func (s *State) GetEpoch(number shared.EpochNumber) *Epoch {
	e, err := s.db.GetEpoch(number)
	if err != nil {
		return nil
	}

	// epoch not found, create new
	if e == nil {
		epochSeed, err := crypto.MixSeed(s.seed, number)
		if err != nil {
			return nil
		}

		e = NewEpochInstance(number, epochSeed)
		err = s.SaveEpoch(e)
		if err != nil {
			return nil
		}
	}

	return e
}

func (s *State) GetPool(poolId shared.PoolId) *Pool {
	return s.Pools[poolId]
}

func (s *State) SavePool(pool *Pool) {
	s.Pools[pool.Id] = pool
}