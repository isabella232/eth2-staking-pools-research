package state

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/herumi/bls-eth-go-binary/bls"
)

type Pool struct {
	id              uint64         // id
	active 			bool
	pubKey          *bls.PublicKey // eth2 validation pubkey
	sortedExecutors []uint64       // ids of the block producers which are executors on this pool
}

func NewPool(
	id              uint64,
	active 			bool,
	pubKey          *bls.PublicKey,
	sortedExecutors []uint64,
	) *Pool {
	return &Pool{
		id:              id,
		active:          active,
		pubKey:          pubKey,
		sortedExecutors: sortedExecutors,
	}
}

func (pool *Pool) Copy() (core.IPool, error) {
	pk := &bls.PublicKey{}
	err := pk.Deserialize(pool.pubKey.Serialize())
	if err != nil {
		return nil, err
	}

	return &Pool{
		id:              pool.id,
		pubKey:          pk,
		sortedExecutors: pool.sortedExecutors,
	}, nil
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

func (pool *Pool) GetPubKey() *bls.PublicKey {
	return pool.pubKey
}

func (pool *Pool) GetSortedExecutors() []uint64 {
	return pool.sortedExecutors
}

func (pool *Pool) SetSortedExecutors(executors []uint64) {
	pool.sortedExecutors = executors
}