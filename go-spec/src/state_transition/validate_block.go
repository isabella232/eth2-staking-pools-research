package state_transition

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-ssz"
)

func (st *StateTransition) PreApplyValidateBlock(state *core.State, header *core.BlockHeader, body *core.BlockBody) error {
	// check necessary vars are not nil
	if len(body.Randao) != 32 {
		return fmt.Errorf("RANDAO should be 32 byte")
	}
	if len(body.ParentBlockRoot) != 32 {
		return fmt.Errorf("parent block root should be 32 byte")
	}

	// validate parent block root
	// Rule 1: need to point to an existing parent block root
	// Rule 2: need to have a higher epoch
	// TODO - block 0?
	foundParent := false
	for _, parent := range state.BlockRoots {
		if bytes.Compare(parent.Bytes, body.ParentBlockRoot) == 0 {
			foundParent = true
			if parent.Epoch >= body.Epoch {
				return fmt.Errorf("new block's parent block root can't be of a future epoch")
			}
		}
	}
	if !foundParent {
		return fmt.Errorf("parent block root not found")
	}

	// verify proposer is expected proposer
	expectedProposer, err := core.GetBlockProposer(state, body.Epoch)
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

	// TODO - validate RANDAO

	// TODO - validate block?

	return nil
}

func (st *StateTransition) PostApplyValidateBlock(newState *core.State, header *core.BlockHeader, body *core.BlockBody) error {
	root := core.GetStateRoot(newState, newState.CurrentEpoch)
	if len(root) == 0 {
		return fmt.Errorf("could not find statet root for epoch %d", newState.CurrentEpoch)
	}

	// validate state root is equal to block
	if hex.EncodeToString(header.StateRoot) != hex.EncodeToString(root[:]) {
		return fmt.Errorf("new block state root is wrong")
	}
	return nil
}