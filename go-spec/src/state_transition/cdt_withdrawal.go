package state_transition

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
)

func (st *StateTransition) validateCDTWithdrawal(state *core.State, withdrawals *core.CDTWithdrawalRequest) error {
	// TODO - check balance sufficient
	// TODO - check approval block
	// TODO - can he withdraw all? should he leave some balance for penalties and so on?
	return nil
}

func (st *StateTransition) ProcessCDTWithdrawals(state *core.State, withdrawals []*core.CDTWithdrawalRequest) error {


	for _, withdrawal := range withdrawals {
		bp := core.GetBlockProducerPubKey(state, withdrawal.PubKey)
		if bp == nil {
			return fmt.Errorf("block producer with pubKey %s not found", hex.EncodeToString(withdrawal.PubKey))
		}

		switch withdrawal.Status {
		case 0: // requested
			if err := st.validateCDTWithdrawal(state, withdrawal); err != nil {
				return err
			}
			 // TODO  - sign the request and broadcast it

			 // to prevent a situation where the request was signed but later reducing the balance fails
			if err := core.DecreaseBPBalance(bp, withdrawal.Amount); err != nil {
				return err
			}
		case 1: // approved

		case 2: // rejected
			// TODO - better define how the rejected status is assigned.
			// we've reduced the balance during signature, increase it if failed
			core.IncreaseBPBalance(bp, withdrawal.Amount)
		}
	}

	return nil
}
