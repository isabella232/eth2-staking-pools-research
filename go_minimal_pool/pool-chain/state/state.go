package state

type DB interface {
	GetEpoch(number uint32) (*Epoch,error)
}

type State struct {
	db DB
	currentEpoch uint32
	pools []*Pool
}

func NewInMemoryState() *State {
	return & State{
		db:           &InMemStateDb{},
		currentEpoch: 0,
		pools:        make([]*Pool, 0),
	}
}

func (s *State) GetEpoch(number uint32) *Epoch {
	e, err := s.db.GetEpoch(number)
	if err != nil {
		return nil
	}

	return e
}

func (s *State) GetCurrentEpoch() *Epoch {
	return s.GetEpoch(s.currentEpoch)
}