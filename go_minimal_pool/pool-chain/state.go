package pool_chain

type DB interface {

}

type State struct {
	db *DB
	currentEpoch *Epoch
	pools []*Pool
}
