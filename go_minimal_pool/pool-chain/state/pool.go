package state

import "github.com/herumi/bls-eth-go-binary/bls"

type Pool struct {
	Id uint8
	Pk *bls.PublicKey
}

func NewPool(id uint8, pk *bls.PublicKey) *Pool {
	return &Pool{
		Id: id,
		Pk: pk,
	}
}