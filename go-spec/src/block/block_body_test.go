package block

import (
	"encoding/hex"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/stretchr/testify/require"
	"testing"
)

func toByte(str string) []byte {
	ret, _ := hex.DecodeString(str)
	return ret
}

func TestBlockBodySSZ(t *testing.T) {
	tests := []struct{
		testName string
		body *BlockBody
		expected []byte
	}{
		{
			testName: "full SSZ",
			body: NewBlockBody(
					12,
					5,
					shared.SliceToByte32(toByte("state root state root state root state root state root")),
					[]*PoolExecutionSummary{
						NewExecutionSummary(
							12,
							5,
							[]*BeaconDuty{
								NewBeaconDuty(
										0,
										12,
										342,
										true,
										[16]byte{1,3,88,12,43,12,89,35,1,0,99,16,63,13,33,0},
									),
								NewBeaconDuty(
									1,
									0,
									343,
									true,
									[16]byte{},
								),
							},
						),
					},
					[]*CreatePoolRequest{
						NewCreatePoolRequest(
							3,
							0,
							5,
							6,
							15,
							toByte("public key"),
							[16]byte{43,12,89,35,99,16,63,13,33,0,1,3,88,12,43,1},
							),
					},
					toByte("parent block root parent block root parent block root parent block root"),
				),
			expected:toByte("d84c6a62eca5740b0b50a350d04d81f5817db9512fb083b2a6a5b1bf69ab75d0"),
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			root,err := test.body.Root()
			require.NoError(t, err)
			require.EqualValues(t, test.expected, root)
		})
	}
}
