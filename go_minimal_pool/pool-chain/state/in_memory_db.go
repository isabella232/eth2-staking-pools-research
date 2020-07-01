package state

type InMemStateDb struct {

}

func (db *InMemStateDb) GetEpoch(number uint32) (*Epoch,error) {
	return &Epoch{Number:number}, nil
}
