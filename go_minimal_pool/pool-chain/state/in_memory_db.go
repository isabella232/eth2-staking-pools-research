package state

import "github.com/bloxapp/eth2-staking-pools-research/minimal_pool/shared"

type InMemStateDb struct {
	epochs map[shared.EpochNumber]*Epoch
}

func NewInMemoryDb() *InMemStateDb {
	return &InMemStateDb{
		epochs: make(map[shared.EpochNumber]*Epoch),
	}
}

func (db *InMemStateDb) SaveEpoch(epoch *Epoch) error {
	db.epochs[epoch.Number] = epoch
	return nil
}

func (db *InMemStateDb) GetEpoch(number shared.EpochNumber) (*Epoch,error) {
	if val, ok := db.epochs[number]; ok {
		return val, nil
	}

	return nil, nil
}
