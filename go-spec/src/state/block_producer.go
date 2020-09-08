package state

import (
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
)

type BlockProducer struct {
	Id				uint64
	PubKey			*bls.PublicKey
	Balance			uint64 // balance on the pool chain (rewards earned)
	Stake			uint64 // stake
	Slashed			bool
	Active 			bool
}

func (bp *BlockProducer) Copy() (*BlockProducer, error) {
	pk := &bls.PublicKey{}
	err := pk.Deserialize(bp.PubKey.Serialize())
	if err != nil {
		return nil, err
	}

	return &BlockProducer{
		Id:      bp.Id,
		PubKey:  pk,
		Balance: bp.Balance,
		Stake:   bp.Stake,
		Slashed: bp.Slashed,
		Active:  bp.Active,
	}, nil
}

func (bp *BlockProducer) IncreaseBalance(change uint64) (newBalance uint64, error error) {
	bp.Balance += change
	return bp.Balance, nil
}

func (bp *BlockProducer) DecreaseBalance(change uint64) (newBalance uint64, error error) {
	if bp.Balance < change {
		return 0, fmt.Errorf("BP %d dosen't have enonugh balance (%d) to decrease (%d)", bp.Id, bp.Balance, change)
	}

	bp.Balance -= change
	return bp.Balance, nil
}
