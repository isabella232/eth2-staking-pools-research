package state_transition

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared/params"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStateCopying(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)

	preRoot, err := ssz.HashTreeRoot(state)
	require.NoError(t, err)

	newState := shared.CopyState(state)
	require.NoError(t, err)

	// test new state and old state ssz
	newStateRoot, err := ssz.HashTreeRoot(newState)
	require.NoError(t, err)
	require.EqualValues(t, preRoot, newStateRoot)

	// test manipulating prams on new state copying
	bp := shared.GetBlockProducer(newState, 0)
	bp.CDTBalance = 100000
	require.NotEqualValues(t, shared.GetBlockProducer(state, 0).CDTBalance, 100000)

	shared.GetPool(newState, 0).Active = false
	require.NotEqualValues(t, shared.GetPool(state, 0).Active, shared.GetPool(newState, 0).Active)

	// test old state root not changed
	postRoot, err := ssz.HashTreeRoot(state)
	require.NoError(t, err)
	require.EqualValues(t, preRoot, postRoot)
}

func TestBlockApplyConsistency(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	block := &core.PoolBlock{
		Slot:                 33,
		Proposer:             2713,
		ParentRoot:           toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		StateRoot:            nil,
		Body:                 &core.PoolBlockBody{
			RandaoReveal:          toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
			NewPoolReq:      nil,
			Attestations:         generateAttestations(
				state,
				128,
				32,
				&core.Checkpoint{Epoch: 0, Root: params.ChainConfig.ZeroHash},
				&core.Checkpoint{Epoch: 1, Root: []byte{}},
				0,
				true,
				0, /* attestation */
			),
		},
	}

	// sign
	sig, err := shared.SignBlock(
		block,
		[]byte(fmt.Sprintf("%d", 2713)),
		[]byte("domain"))
	require.NoError(t, err)
	signed := &core.SignedPoolBlock{
		Block:                block,
		Signature:            sig.Serialize(),
	}


	preRoot,err := ssz.HashTreeRoot(state)
	require.NoError(t, err)

	var postRoot []byte
	for i := 0 ; i < 10 ; i++ {
		newState := shared.CopyState(state)
		st := NewStateTransition()
		err := st.ProcessBlock(newState, signed)
		require.NoError(t, err)

		if i != 0 {
			require.EqualValues(t, postRoot, shared.GetStateRoot(newState, 1))
		}

		postRoot = shared.GetStateRoot(newState, 1)
	}

	post,err := ssz.HashTreeRoot(state)
	require.NoError(t, err)
	require.EqualValues(t, preRoot, post)
}