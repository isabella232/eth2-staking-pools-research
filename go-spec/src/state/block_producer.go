package state

import (
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
)

type BlockProducer struct {
	id      uint64
	pubKey  []byte
	balance uint64 // balance on the pool chain (rewards earned)
	stake   uint64 // stake
	slashed bool
	active  bool
	exitEpoch uint64
}

func NewBlockProducer(
	id      uint64,
	pubKey  []byte,
	balance uint64,
	stake   uint64,
	slashed bool,
	active  bool,
	) *BlockProducer {
	return &BlockProducer{
		id:      id,
		pubKey:  pubKey,
		balance: balance,
		stake:   stake,
		slashed: slashed,
		active:  active,
	}
}

func (bp *BlockProducer) GetId() uint64 {
	return bp.id
}

func (bp *BlockProducer) GetPubKey() (*bls.PublicKey, error) {
	ret := &bls.PublicKey{}
	err := ret.Deserialize(bp.pubKey)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (bp *BlockProducer) SetPubKey(pk []byte) {
	bp.pubKey = pk
}

func (bp *BlockProducer) GetBalance() uint64 {
	return bp.balance
}

func (bp *BlockProducer) GetStake() uint64 {
	return bp.stake
}

func (bp *BlockProducer) IsSlashed() bool {
	return bp.slashed
}

func (bp *BlockProducer) IsActive() bool {
	return bp.active
}

func (bp *BlockProducer) SetExited(atEpoch uint64) {
	bp.active = false
	bp.exitEpoch = atEpoch
}

func (bp *BlockProducer) ExitEpoch() uint64 {
	if bp.IsActive() {
		return 0
	}
	return bp.exitEpoch
}

func (bp *BlockProducer) IncreaseBalance(change uint64) (newBalance uint64, error error) {
	bp.balance += change
	return bp.balance, nil
}

func (bp *BlockProducer) DecreaseBalance(change uint64) (newBalance uint64, error error) {
	if bp.balance < change {
		return 0, fmt.Errorf("BP %d dosen't have enonugh balance (%d) to decrease (%d)", bp.id, bp.balance, change)
	}

	bp.balance -= change
	return bp.balance, nil
}
