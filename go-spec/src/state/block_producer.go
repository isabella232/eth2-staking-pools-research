package state

import (
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
)

type BlockProducer struct {
	Id        uint64
	PubKey    []byte
	Balance   uint64 // Balance on the pool chain (rewards earned)
	Stake     uint64 // Stake
	Slashed   bool
	Active    bool
	ExitEpoch uint64
}

func (bp *BlockProducer) GetId() uint64 {
	return bp.Id
}

func (bp *BlockProducer) GetPubKey() (*bls.PublicKey, error) {
	ret := &bls.PublicKey{}
	err := ret.Deserialize(bp.PubKey)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (bp *BlockProducer) SetPubKey(pk []byte) {
	bp.PubKey = pk
}

func (bp *BlockProducer) GetBalance() uint64 {
	return bp.Balance
}

func (bp *BlockProducer) GetStake() uint64 {
	return bp.Stake
}

func (bp *BlockProducer) IsSlashed() bool {
	return bp.Slashed
}

func (bp *BlockProducer) IsActive() bool {
	return bp.Active
}

func (bp *BlockProducer) SetExited(atEpoch uint64) {
	bp.Active = false
	bp.ExitEpoch = atEpoch
}

func (bp *BlockProducer) GetExitEpoch() uint64 {
	if bp.IsActive() {
		return 0
	}
	return bp.ExitEpoch
}

func (bp *BlockProducer) IncreaseBalance(change uint64) (newBalance uint64, error error) {
	bp.Balance += change
	return bp.Balance, nil
}

func (bp *BlockProducer) DecreaseBalance(change uint64) (newBalance uint64, error error) {
	if bp.Balance < change {
		return 0, fmt.Errorf("BP %d dosen't have enonugh Balance (%d) to decrease (%d)", bp.Id, bp.Balance, change)
	}

	bp.Balance -= change
	return bp.Balance, nil
}
