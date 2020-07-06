package state

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

func getSeed(str string) [32]byte {
	var ret [32]byte
	_seed, _ := hex.DecodeString(str)
	copy(ret[:],_seed)

	return ret
}

func TestShufflePools(t *testing.T) {
	tests := []struct{
		testName string
		seed [32]byte
		input []uint32
		cntPools uint8
		poolSize uint8
		roundCount uint8
		expectedShufflee map[uint8][]uint32
	}{
		{
			testName: "1 pool of 3",
			seed: getSeed("b581262ce281d1e9deaf2f0158d7cd05217f1196d95956c5f55d837ccc3c8a9"),
			input: []uint32{
				1,2,3,
			},
			cntPools: 1,
			poolSize:3,
			roundCount: 10,
			expectedShufflee: map[uint8][]uint32{
				1: []uint32{2,1,3},
			},
		},
		{
			testName: "2 pool of 3",
			seed: getSeed("b581262ce281d1e9deaf2f0158d7cd05217f1196d95956c5f55d837ccc3c8a9"),
			input: []uint32{
				1,2,3,4,5,6,
			},
			cntPools: 2,
			poolSize:3,
			roundCount: 10,
			expectedShufflee: map[uint8][]uint32{
				1: []uint32{6,1,4},
				2: []uint32{2,3,5},
			},
		},
		{
			testName: "3 pool of 3",
			seed: getSeed("b581262ce281d1e9deaf2f0158d7cd05217f1196d95956c5f55d837ccc3c8a9"),
			input: []uint32{
				1,2,3,4,5,6,7,8,9,
			},
			cntPools: 3,
			poolSize:3,
			roundCount: 10,
			expectedShufflee: map[uint8][]uint32{
				1: []uint32{7,3,2},
				2: []uint32{8,1,9},
				3: []uint32{5,4,6},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			res,err := shufflePools(
					test.input,
					test.seed,
					test.roundCount,
					test.cntPools,
					test.poolSize,
				)
			require.NoError(t, err)
			require.Equal(t, test.expectedShufflee, res)
		})
	}
}