package state_transition

import (
	"bytes"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-ssz"
)

func (st *StateTransition) ValidateBlock(state *core.State, header *core.BlockHeader, body *core.BlockBody) error {
	// verify proposer is expected proposer
	expectedProposer, err := core.GetBlockProposer(state, state.CurrentEpoch)
	if err != nil {
		return err
	}
	proposerId := body.GetProposer()
	if expectedProposer != proposerId {
		return fmt.Errorf("block expectedProposer is worng, expected %d but received %d", expectedProposer, proposerId)
	}

	// verify header block root matches
	exectedRoot,err := ssz.HashTreeRoot(body)
	if err != nil {
		return err
	}
	if bytes.Compare(exectedRoot[:], header.GetBlockRoot()) != 0 {
		return fmt.Errorf("signed block root does not match body root")
	}

	// validate signature
	proposer := core.GetBlockProducer(state, proposerId)
	if proposer == nil {
		return fmt.Errorf("proposer not found")
	}
	sig := &bls.Sign{}
	err = sig.Deserialize(header.Signature)
	if err != nil {
		return err
	}
	pk := &bls.PublicKey{}
	err = pk.Deserialize(proposer.GetPubKey())
	if err != nil {
		return err
	}
	if res := sig.VerifyHash(pk, header.BlockRoot); !res {
		return fmt.Errorf("signature did not verify")
	}

	// TODO - validate block?

	return nil
}
