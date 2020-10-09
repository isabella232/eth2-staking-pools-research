package state_transition

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/prysmaticlabs/go-ssz"
)

func (st *StateTransition) ProcessBlock(state *core.State, signedBlock *core.SignedPoolBlock) error {
	if err := processBlockHeader(state, signedBlock); err != nil {
		return err
	}
	if err := processRANDAO(state, signedBlock.Block); err != nil {
		return err
	}
	// operations
	if err := st.processBlockAttestations(state, signedBlock.Block.Body.Attestations); err != nil {
		return err
	}
	if err := st.processNewPoolRequests(state, signedBlock.Block.Body.NewPoolReq); err != nil {
		return err
	}
	return nil
}

func (st *StateTransition) processBlockForStateRoot(state *core.State, signedBlock *core.SignedPoolBlock) error {
	if err := processBlockHeaderNoVerify(state, signedBlock); err != nil {
		return err
	}
	if err := processRANDAONoVerify(state, signedBlock.Block); err != nil {
		return err
	}
	for _, att := range signedBlock.Block.Body.Attestations {
		if err := processAttestationNoSigVerify(st, state, att); err != nil {
			return err
		}
	}
	if err := st.processNewPoolRequests(state, signedBlock.Block.Body.NewPoolReq); err != nil {
		return err
	}
	return nil
}

// ProcessBlockHeader validates a block by its header.
//
// Spec pseudocode definition:
//
//  def process_block_header(state: BeaconState, block: BeaconBlock) -> None:
//    # Verify that the slots match
//    assert block.slot == state.slot
//     # Verify that proposer index is the correct index
//    assert block.proposer_index == get_beacon_proposer_index(state)
//    # Verify that the parent matches
//    assert block.parent_root == hash_tree_root(state.latest_block_header)
//    # Save current block as the new latest block
//    state.latest_block_header = BeaconBlockHeader(
//        slot=block.slot,
//        parent_root=block.parent_root,
//        # state_root: zeroed, overwritten in the next `process_slot` call
//        body_root=hash_tree_root(block.block),
//		  # signature is always zeroed
//    )
//    # Verify proposer is not slashed
//    proposer = state.validators[get_beacon_proposer_index(state)]
//    assert not proposer.slashed
//    # Verify proposer signature
//    assert bls_verify(proposer.pubkey, signing_root(block), block.signature, get_domain(state, DOMAIN_BEACON_PROPOSER))
func processBlockHeaderNoVerify(state *core.State, signedBlock *core.SignedPoolBlock) error {
	block := signedBlock.Block

	// slot
	if state.CurrentSlot != block.Slot {
		return fmt.Errorf("block slot doesn't match state slot")
	}

	// proposer
	expectedProposer, err := shared.BlockProposer(state, block.Slot)
	if err != nil {
		return err
	}
	proposerId :=  block.GetProposer()
	if expectedProposer != proposerId {
		return fmt.Errorf("block expectedProposer is worng, expected %d but received %d", expectedProposer, proposerId)
	}

	// parent
	root,err := ssz.HashTreeRoot(state.LatestBlockHeader)
	if err != nil {
		return err
	}
	if !bytes.Equal(block.ParentRoot, root[:]) {
		return fmt.Errorf("parent block root doesn't match, expected %s", hex.EncodeToString(root[:]))
	}

	// save
	root,err = ssz.HashTreeRoot(block.Body)
	if err != nil {
		return err
	}
	state.LatestBlockHeader = &core.PoolBlockHeader{
		Slot:                 block.Slot,
		ProposerIndex:        block.Proposer,
		ParentRoot:           block.ParentRoot,
		BodyRoot:             root[:],
	}

	// TODO - verify proposer is not slashed

	return nil
}

func processBlockHeader(state *core.State, signedBlock *core.SignedPoolBlock) error {
	if err := processBlockHeaderNoVerify(state, signedBlock); err != nil {
		return err
	}
	if err := verifyBlockSig(state, signedBlock); err != nil {
		return err
	}
	return nil
}

func verifyBlockSig(state *core.State, signedBlock *core.SignedPoolBlock) error {
	block := signedBlock.Block

	// verify sig
	proposer := shared.GetBlockProducer(state, block.GetProposer())
	if proposer == nil {
		return fmt.Errorf("proposer not found")
	}
	if err := shared.VerifyBlockSigningRoot(block, proposer.GetPubKey(), signedBlock.Signature, []byte("domain")); err != nil { // TODO - domain not hard coded
		return err
	}
	return nil
}

func processRANDAO (state *core.State, body *core.PoolBlock) error {
	return nil
}

func processRANDAONoVerify(state *core.State, body *core.PoolBlock) error {
	return nil
}