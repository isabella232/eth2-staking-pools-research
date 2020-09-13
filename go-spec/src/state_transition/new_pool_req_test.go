package state_transition

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

func TestCreatedNewPoolReq(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	_, body := GenerateCreatePoolHeadAndBody(state)

	st := NewStateTransition()

	newState, err := st.ApplyBlock(state, body)
	require.NoError(t, err)

	// check created
	require.Equal(t, 6, len(newState.Pools))

	// check rewards
	participation := []byte{43,12,89,35,99,16,63,13,33,0,1,3,88,12,43,1}
	committee, err := core.PoolCommittee(newState, 6, 1)
	sort.Slice(committee, func(i int, j int) bool {
		return committee[i] < committee[j]
	})
	require.NoError(t, err)

	// test penalties/ rewards
	for i := uint64(0) ; i < core.TestConfig().DKGParticipantsNumber ; i++ {
		bp := core.GetBlockProducer(newState, committee[i])
		if shared.IsBitSet(participation, i) {
			require.EqualValues(t, 2000, bp.Balance)
		} else {
			require.EqualValues(t, 0, bp.Balance)
		}
	}

	// leader reward
	bp := core.GetBlockProducer(newState, 1)
	require.EqualValues(t, 4000, bp.Balance)

	// pool data
	pool := core.GetPool(newState, 6)
	require.NotNil(t, pool)
	require.EqualValues(t, toByte("a3b9110ec26cbb02e6182fab4dcb578d17411f26e41f16aad99cfce51e9bc76ce5e7de00a831bbcadd1d7bc0235c945d"), pool.PubKey)
	require.EqualValues(t, committee, pool.SortedCommittee)
}

func TestNotCreatedNewPoolReq(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	_, body := GenerateNotCreatePoolHeadAndBody(state)

	st := NewStateTransition()

	newState, err := st.ApplyBlock(state, body)
	require.NoError(t, err)

	// check not created
	require.Equal(t, 5, len(newState.Pools))

	// check penalties
	committee, err := core.PoolCommittee(newState, 6, 1)
	sort.Slice(committee, func(i int, j int) bool {
		return committee[i] < committee[j]
	})
	require.NoError(t, err)

	// test penalties/ rewards
	for i := uint64(0) ; i < core.TestConfig().DKGParticipantsNumber ; i++ {
		bp := core.GetBlockProducer(newState, committee[i])
		require.EqualValues(t, 0, bp.Balance)
	}

	// leader reward
	bp := core.GetBlockProducer(newState, 1)
	require.EqualValues(t, 1000, bp.Balance)
}

func TestCreatedNewPoolReqWithExistingId(t *testing.T) {
	t.Skipf("skipping as we can't build state root becasuse state apply errors")

	//require.NoError(t, bls.Init(bls.BLS12_381))
	//require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))
	//
	//state := generateTestState(t)
	//_, body := GenerateCreatePoolWithExistingIdHeadAndBody(state)
	//
	//st := NewStateTransition()
	//
	//_, err := st.ApplyBlock(state, body)
	//require.EqualError(t, err, "new pool id == req id, this is already exists")
}