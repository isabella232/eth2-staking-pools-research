package state_transition

import "github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"

func (st *StateTransition) ProcessBlockAttestations(state *core.State, attestations []*core.Attestation) error {
	for _, att := range attestations {
		if err := st.processAttestation(state, att); err != nil {
			return err
		}
	}
	return nil
}

func (st *StateTransition) processAttestation(state *core.State, attestation *core.Attestation) error {
	// TODO - valdiate signatures
	// TODO - valdiate 2/3 sign
	// TODO - validate epoch, slot, inclusion distance

	// process execution summaries
	if err := st.processExecutionSummaries(state, attestation.Data.ExecutionSummaries); err != nil {
		return err
	}

	return nil
}
