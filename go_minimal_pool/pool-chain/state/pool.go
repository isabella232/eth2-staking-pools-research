package state

import (
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
)

type Pool struct {
	Id shared.PoolId
	Size shared.PoolSize
	Pk *bls.PublicKey
}

func NewPool(id shared.PoolId, size shared.PoolSize, pk *bls.PublicKey) *Pool {
	return &Pool{
		Id: id,
		Size:size,
		Pk: pk,
	}
}