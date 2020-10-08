package state_transition

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProcessBlockHeader(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))
	tests := []struct{
		name              string
		block             *core.PoolBlock
		proposerId        uint64
		expectedError     error
		useCorretBodyRoot bool
	}{
		{
			name: "valid sig",
			block: &core.PoolBlock{
				Proposer:        2713,
				Slot:            35,
				Body: &core.PoolBlockBody{
					RandaoReveal:         toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
				},
				ParentRoot: toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
			},
			proposerId: 2713,
			expectedError: nil,
			useCorretBodyRoot: true,
		},
		{
			name: "invalid sig",
			block: &core.PoolBlock{
				Proposer:        2713,
				Slot:            35,
				Body: &core.PoolBlockBody{
					RandaoReveal:         toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
				},
				ParentRoot: toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
			},
			proposerId: 2712,
			expectedError: fmt.Errorf("signature did not verify"),
			useCorretBodyRoot: true,
		},
		{
			name: "wrong proposer",
			block: &core.PoolBlock{
				Proposer:        2,
				Slot:            35,
				Body: &core.PoolBlockBody{
					RandaoReveal:         toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
				},
				ParentRoot: toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
			},
			proposerId: 2,
			expectedError: fmt.Errorf("block expectedProposer is worng, expected 2713 but received 2"),
			useCorretBodyRoot: true,
		},
		{
			name: "invalid proposer",
			block: &core.PoolBlock{
				Proposer:        4550000000,
				Slot:            35,
				Body: &core.PoolBlockBody{
					RandaoReveal:         toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
				},
				ParentRoot: toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
			},
			proposerId: 2,
			expectedError: fmt.Errorf("block expectedProposer is worng, expected 2713 but received 4550000000"),
			useCorretBodyRoot: true,
		},
		{
			name: "invalid block root",
			block: &core.PoolBlock{
				Proposer:        2713,
				Slot:            35,
				Body: &core.PoolBlockBody{
					RandaoReveal:         toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
				},
				ParentRoot: toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
			},
			proposerId: 2713,
			expectedError: fmt.Errorf("signed block root does not match block root"),
			useCorretBodyRoot: false,
		},
		{
			name: "RANDAO too small",
			block: &core.PoolBlock{
				Proposer:        2713,
				Slot:            35,
				Body: &core.PoolBlockBody{
					RandaoReveal:         toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6"),
				},
				ParentRoot: toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
			},
			proposerId: 2713,
			expectedError: fmt.Errorf("RANDAO should be 32 byte"),
			useCorretBodyRoot: true,
		},
		{
			name: "RANDAO too big",
			block: &core.PoolBlock{
				Proposer:        2713,
				Slot:            35,
				Body: &core.PoolBlockBody{
					RandaoReveal:         toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6ddd"),
				},
				ParentRoot: toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
			},
			proposerId: 2713,
			expectedError: fmt.Errorf("RANDAO should be 32 byte"),
			useCorretBodyRoot: true,
		},
		{
			name: "invalid parent block root",
			block: &core.PoolBlock{
				Proposer:        2713,
				Slot:            35,
				Body: &core.PoolBlockBody{
					RandaoReveal:         toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
				},
				ParentRoot: toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e42"),
			},
			proposerId: 2713,
			expectedError: fmt.Errorf("parent block root not found"),
			useCorretBodyRoot: true,
		},
		{
			name: "block from the past",
			block: &core.PoolBlock{
				Proposer:        2713,
				Slot:            35,
				Body: &core.PoolBlockBody{
					RandaoReveal:         toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
				},
				ParentRoot: toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
			},
			proposerId: 2713,
			expectedError: fmt.Errorf("new block's parent block root can't be of a future epoch"),
			useCorretBodyRoot: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := generateTestState(t, 3)

			// block root
			root,err := ssz.HashTreeRoot(test.block)
			require.NoError(t, err)

			// sign
			sk := []byte(fmt.Sprintf("%d", test.proposerId))
			sig, err := shared.SignBlock(test.block, sk, []byte("domain")) // TODO - dynamic domain
			require.NoError(t, err)

			if !test.useCorretBodyRoot {
				root[0] = 0
			}

			// header
			signed := &core.SignedPoolBlock{
				Block:                test.block,
				Signature:            sig.Serialize(),
			}

			if test.expectedError != nil {
				require.EqualError(t, processBlockHeader(state, signed), test.expectedError.Error())
			} else {
				require.NoError(t, processBlockHeader(state, signed))
			}
		})
	}
}