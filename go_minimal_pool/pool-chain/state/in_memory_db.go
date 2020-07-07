package state

type InMemStateDb struct {
	epochs map[uint32]*Epoch
}

func NewInMemoryDb() *InMemStateDb {
	return &InMemStateDb{
		epochs: make(map[uint32]*Epoch),
	}
}

func (db *InMemStateDb) SaveEpoch(epoch *Epoch) error {
	db.epochs[epoch.Number] = epoch
	return nil
}

func (db *InMemStateDb) GetEpoch(number uint32) (*Epoch,error) {
	if val, ok := db.epochs[number]; ok {
		return val, nil
	}

	return nil, nil
}
