package state

import (
	"github.com/herumi/bls-eth-go-binary/bls"
)

type Pool struct {
	Id              uint64 // Id
	Active          bool
	PubKey          []byte   // eth2 validation pubkey
	SortedExecutors []uint64 // ids of the block producers which are executors on this pool
}

func (pool *Pool) IsActive() bool {
	return pool.Active
}

func (pool *Pool) SetActive(status bool) {
	pool.Active = status
}

func (pool *Pool) GetId() uint64 {
	return pool.Id
}

func (pool *Pool) GetPubKey() (*bls.PublicKey, error) {
	ret := &bls.PublicKey{}
	err := ret.Deserialize(pool.PubKey)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (pool *Pool) GetSortedExecutors() []uint64 {
	return pool.SortedExecutors
}

func (pool *Pool) SetSortedExecutors(executors []uint64) {
	pool.SortedExecutors = executors
}