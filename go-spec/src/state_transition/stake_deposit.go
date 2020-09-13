package state_transition

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
)

func (st *StateTransition) validateStakeDeposits(state *core.State, deposits []*core.StakeDeposit) error {
	// TODO - eth2 (prysm) validates deposits using merkle proofs, we should consider using them
	return nil
}

func (st *StateTransition) ProcessStakeDeposits(state *core.State, deposits []*core.StakeDeposit) error {
	if err := st.validateStakeDeposits(state, deposits); err != nil {
		return err
	}

	for _, deposit := range deposits {
		bp := core.GetBlockProducerPubKey(state, deposit.PubKey)
		if bp == nil {
			return fmt.Errorf("block producer with pubKey %s not found", hex.EncodeToString(deposit.PubKey))
		}

		bp.Stake += deposit.Amount
	}

	return nil
}
