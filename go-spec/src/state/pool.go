package state

import (
	"github.com/herumi/bls-eth-go-binary/bls"
)

type Pool struct {
	id              uint64         // id
	active 			bool
	pubKey          []byte // eth2 validation pubkey
	sortedExecutors []uint64       // ids of the block producers which are executors on this pool
}

func NewPool(
	id              uint64,
	active 			bool,
	pubKey          []byte,
	sortedExecutors []uint64,
	) *Pool {
	return &Pool{
		id:              id,
		active:          active,
		pubKey:          pubKey,
		sortedExecutors: sortedExecutors,
	}
}

func (pool *Pool) IsActive() bool {
	return pool.active
}

func (pool *Pool) SetActive(status bool) {
	pool.active = status
}

func (pool *Pool) GetId() uint64 {
	return pool.id
}

func (pool *Pool) GetPubKey() (*bls.PublicKey, error) {
	ret := &bls.PublicKey{}
	err := ret.Deserialize(pool.pubKey)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (pool *Pool) GetSortedExecutors() []uint64 {
	return pool.sortedExecutors
}

func (pool *Pool) SetSortedExecutors(executors []uint64) {
	pool.sortedExecutors = executors
}