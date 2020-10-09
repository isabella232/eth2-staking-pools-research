package state_transition

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandaoSeedMix(t *testing.T) {
	t.Skipf("randao not implementd yet")
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t, 3)
	block := &core.PoolBlock{
		Proposer:        2713,
		Slot:            35,
		Body: &core.PoolBlockBody{
			RandaoReveal:         toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		},
		ParentRoot: toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
	}

	require.NoError(t, processRANDAO(state, block))

	newsSeed, err := shared.GetEpochSeed(state, 35)
	require.NoError(t, err)
	require.EqualValues(t, toByte("e4a17401658219365021cf584f4758d4b22ec861d9653e8249c8b4f73285a909"), newsSeed)
}


func TestBlockPostValidation(t *testing.T) {
	tests := []struct{
		name          string
		block         *core.PoolBlock
		expectedError error
	}{
		{
			name: "valid post state root",
			block: &core.PoolBlock{
				Proposer:        27,
				Slot:            4,
				Body: &core.PoolBlockBody{
					RandaoReveal:         toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
				},
				ParentRoot: toByte("ddfc878b39d0e6807c673b2775788ebba122c37786bea8b23a407c532ce309fe"),
				StateRoot: toByte("70069e2e57bf907263fa1edf3e272b65a5652a89a388c35586d70e28fcc19977"),
			},
			expectedError: nil,
		},
		{
			name: "invalid post state root",
			block: &core.PoolBlock{
				Proposer:        27,
				Slot:            4,
				Body: &core.PoolBlockBody{
					RandaoReveal:         toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
				},
				ParentRoot: toByte("ddfc878b39d0e6807c673b2775788ebba122c37786bea8b23a407c532ce309fe"),
				StateRoot: toByte("70069e2e57bf907263fa1edf3e272b65a5652a89a388c35586d70e28fcc19976"),
			},
			expectedError: fmt.Errorf("new block state root is wrong, expected 70069e2e57bf907263fa1edf3e272b65a5652a89a388c35586d70e28fcc19977"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := generateTestState(t, 3)

			// sign
			sk := []byte(fmt.Sprintf("%d", test.block.Proposer))
			sig, err := shared.SignBlock(test.block, sk, []byte("domain")) // TODO - dynamic domain
			require.NoError(t, err)

			// header
			signed := &core.SignedPoolBlock{
				Block:                test.block,
				Signature:            sig.Serialize(),
			}

			st := NewStateTransition()
			_, err = st.ExecuteStateTransition(state, signed)
			if test.expectedError != nil {
				require.EqualError(t, err, test.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
