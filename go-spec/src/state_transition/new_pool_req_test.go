package state_transition

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatedNewPoolReq(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	head, body := GenerateCreatePoolHeadAndBody(t)

	st := NewStateTransition()

	newState, err := st.ApplyBlockBody(state, head, body)
	require.NoError(t, err)

	// check created
	require.Equal(t, 6, len(newState.Pools))

	// check rewards
	participation := []byte{43,12,89,35,99,16,63,13,33,0,1,3,88,12,43,1}
	committee, err := core.DKGCommittee(state, 3, 0)
	require.NoError(t, err)

	// test penalties/ rewards
	for i := uint64(0) ; i < core.TestConfig().DKGParticipantsNumber ; i++ {
		bp := core.GetBlockProducer(state, committee[i])
		if shared.IsBitSet(participation, i) {
			require.EqualValues(t, 2000, bp.Balance)
		} else {
			require.EqualValues(t, 0, bp.Balance)
		}
	}

	// leader reward
	bp := core.GetBlockProducer(state, 0)
	require.EqualValues(t, 4000, bp.Balance)
}

func TestNotCreatedNewPoolReq(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	head, body := GenerateNotCreatePoolHeadAndBody(t)

	st := NewStateTransition()

	newState, err := st.ApplyBlockBody(state, head, body)
	require.NoError(t, err)

	// check not created
	require.Equal(t, 5, len(newState.Pools))

	// check penalties
	committee, err := core.DKGCommittee(state, 3, 0)
	require.NoError(t, err)

	// test penalties/ rewards
	for i := uint64(0) ; i < core.TestConfig().DKGParticipantsNumber ; i++ {
		bp := core.GetBlockProducer(state, committee[i])
		require.EqualValues(t, 0, bp.Balance)
	}

	// leader reward
	bp := core.GetBlockProducer(state, 0)
	require.EqualValues(t, 1000, bp.Balance)
}
