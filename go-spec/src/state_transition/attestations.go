package state_transition

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-ssz"
)

func (st *StateTransition) ProcessBlockAttestations(state *core.State, blockBody *core.BlockBody) error {
	attestations := blockBody.Attestations
	for _, att := range attestations {
		if err := st.processAttestation(state, att, blockBody.Slot); err != nil {
			return err
		}
	}
	return nil
}



func (st *StateTransition) processAttestation(state *core.State, attestation *core.Attestation, slot uint64) error {
	// TODO - valdiate signatures
	// TODO - valdiate 2/3 sign
	// TODO - validate epoch, slot, inclusion distance

	if err := validateSignature(state, attestation, slot); err != nil {
		return err
	}

	// process execution summaries
	if err := st.processExecutionSummaries(state, attestation.Data.ExecutionSummaries); err != nil {
		return err
	}

	return nil
}

func validateSignature(state *core.State, attestation *core.Attestation, slot uint64) error {
	// reconstruct committee
	expectedCommittee, err := shared.SlotCommittee(state, slot, uint64(attestation.Data.CommitteeIndex))
	if err != nil {
		return err
	}

	// get pubkeys by aggregation bits
	pks := make([]bls.PublicKey,0)

	for i, id := range expectedCommittee {
		bp := core.GetBlockProducer(state, id)
		if bp == nil {
			return fmt.Errorf("BP %d is inactivee ", id)
		}

		// deserialize pk and aggregate
		if attestation.AggregationBits.BitAt(uint64(i)) {
			pk := bls.PublicKey{}
			err := pk.Deserialize(bp.PubKey)
			if err != nil {
				return err
			}
			pks = append(pks, pk)
		}
	}

	// verify
	sig := &bls.Sign{}
	err = sig.Deserialize(attestation.Signature)
	if err != nil {
		return err
	}
	root, err := ssz.HashTreeRoot(attestation.Data)
	if err != nil {
		return err
	}
	res := sig.FastAggregateVerify(pks, root[:])
	if !res {
		return fmt.Errorf("attestation signature not vrified")
	}

	return nil
}
