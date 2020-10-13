package state_transition

import (
	"bytes"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandaoRevealMix(t *testing.T) {
	state := generateTestState(t, 3)
	proposer, err := shared.BlockProposer(state, 2)
	require.NoError(t, err)
	// get data
	data, domain, err := shared.RANDAOSigningData(state)
	require.NoError(t, err)
	// sign
	sig, err := shared.SignRandao(data, domain, []byte(fmt.Sprintf("%d", proposer)))
	require.NoError(t, err)
	require.NoError(t, processRANDAONoVerify(state,
		&core.PoolBlock{
			Slot:                 2,
			Proposer:             0,
			ParentRoot:           nil,
			StateRoot:            nil,
			Body:                 &core.PoolBlockBody{
				RandaoReveal: sig.Serialize(),
				Attestations: nil,
			},
		}),
	)
	require.True(t, bytes.Equal(state.Seeds[3].Bytes, toByte("5265f3ab14b3a21e3a39002a378c020a298711108830a3f2198d182953593362")))
}
