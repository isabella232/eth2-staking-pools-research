package state

import "github.com/bloxapp/eth2-staking-pools-research/minimal_pool/crypto"

type DB interface {
	// will return nil,nil if epoch not found
	GetEpoch(number uint32) (*Epoch,error)
	SaveEpoch(epoch *Epoch) error
}

type State struct {
	db           DB
	currentEpoch uint32
	Pools        map[uint8]*Pool
	seed         [32]byte
}

func NewInMemoryState(seed [32]byte) *State {
	return & State{
		db:           NewInMemoryDb(),
		currentEpoch: 0,
		Pools:        make(map[uint8]*Pool),
		seed:         seed,
	}
}

func (s *State) SaveEpoch(epoch *Epoch) error {
	return s.db.SaveEpoch(epoch)
}

func (s *State) GetEpoch(number uint32) *Epoch {
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

func (s *State) GetCurrentEpoch() *Epoch {
	return s.GetEpoch(s.currentEpoch)
}

func (s *State) GetPool(poolId uint8) *Pool {
	return s.Pools[poolId]
}

func (s *State) SavePool(pool *Pool) {
	s.Pools[pool.Id] = pool
}