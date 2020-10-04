package shared

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCommitteeStructure(t *testing.T) {
	tests := []struct{
		name string
		bpCount uint64
		expectedSlotsAndCommittees map[uint64][]int
		expectedTotalCommittees uint64
	}{
		{
			name: "not all slots filled, not even committee sizes",
			bpCount: 1000,
			expectedSlotsAndCommittees: map[uint64][]int{
				0: []int{143},
				1: []int{143},
				2: []int{143},
				3: []int{143},
				4: []int{143},
				5: []int{143},
				6: []int{142},
			},
			expectedTotalCommittees: 7,
		},
		{
			name: "not all slots filled, even committee sizes",
			bpCount: 1001,
			expectedSlotsAndCommittees: map[uint64][]int{
				0: []int{143},
				1: []int{143},
				2: []int{143},
				3: []int{143},
				4: []int{143},
				5: []int{143},
				6: []int{143},
			},
			expectedTotalCommittees: 7,
		},
		{
			name: "not all slots filled, even committee sizes",
			bpCount: 1102,
			expectedSlotsAndCommittees: map[uint64][]int{
				0: []int{138},
				1: []int{138},
				2: []int{138},
				3: []int{138},
				4: []int{138},
				5: []int{138},
				6: []int{137},
			},
			expectedTotalCommittees: 8,
		},
		{
			name: "all slots filled, even committee sizes",
			bpCount: 640,
			expectedSlotsAndCommittees: map[uint64][]int{
				0: []int{128},
				1: []int{128},
				2: []int{128},
				3: []int{128},
				4: []int{128},
			},
			expectedTotalCommittees: 5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			activebps := make([]uint64, test.bpCount)
			for i := range activebps {
				activebps[i] = uint64(i)
			}
			committees := CommitteeStructure(activebps)

			for slot, slotCommittee := range committees {
				t.Run(fmt.Sprintf("%s (slot %d)", test.name, slot), func(t *testing.T) {
					require.NotNil(t, slotCommittee)

					// randome committee sample
					if expected := test.expectedSlotsAndCommittees[slot]; expected != nil {
						require.EqualValues(t, len(expected), len(slotCommittee)) // number of committees
						for i, cnt := range expected { // each committee size
							require.EqualValues(t, cnt, len(slotCommittee[i]))
						}
					}
				})
			}

			t.Run(fmt.Sprintf("%s - total committee count", test.name), func(t *testing.T) {
				// total committee size
				allCommittees := 0
				for _, slot := range committees {
					allCommittees += len(slot)
				}
				require.EqualValues(t, test.expectedTotalCommittees, allCommittees)
			})
		})
	}
}

func TestCommitteeShuffling(t *testing.T) {
	// test state
	pools := 128
	bpInPool := 128
	bps := make([]*core.BlockProducer, pools * bpInPool)
	for i := 0 ; i < len(bps) ; i++ {
		bps[i] = &core.BlockProducer{
			Id:      uint64(i),
			Stake:   0,
			Slashed: false,
			Active:  true,
			PubKey:  []byte(fmt.Sprintf("pubkey %d", i)),
		}
	}

	state := &core.State{
		GenesisTime:          0,
		CurrentEpoch:         0,
		BlockRoots:           nil,
		StateRoots:           nil,
		Seeds:                []*core.EpochAndBytes{
			&core.EpochAndBytes{
				Epoch:               0,
				Bytes:                 []byte("seedseedseedseedseedseedseedseed"),
			},
			&core.EpochAndBytes{
				Epoch:               1,
				Bytes:                 []byte("sdddseedseedseedseedseedseedseed"),
			},
		},
		BlockProducers:       bps,
		Pools:                nil,
	}

	tests := []struct{
		name                         string
		epoch                        uint64
		slot                         uint64
		poolId                       uint64
		committeeId                  uint64
		expectedVaultCommittee       []uint64
		expectedAttestationCommittee []uint64
		expectedBlockProposer        uint64
	}{
		{
			name:                         "slot 35, pool id 1, committee id 1",
			epoch:                        1,
			slot:                         35,
			poolId:                       1,
			committeeId:                  1,
			expectedVaultCommittee:       []uint64{2253,666,4476,11882,15304,12496,4906,3489,583,4941,4760,4311,3266,2992,5122,9382,10153,924,9538,6143,10782,7127,11776,7620},
			expectedAttestationCommittee: []uint64{3770,8672,4345,8313,554,13541,9392,15755,9175,3468,13221,14464,9703,1249,10161,7672,2445,2153,4301,909,89,8375,15102,1659,1004,11653,2384,10433,8814,11685,11878,5170,5267,11478,8705,10542,8344,9406,5423,6875,12261,14902,14441,9353,9522,10990,13501,7481,7021,13720,3674,2960,5299,10204,7102,10702,15683,16337,4411,720,3930,1467,5740,7358,14810,12470,11116,13107,6555,7824,9483,10718,4254,8710,6622,14673,2646,13120,14641,3639,15017,10171,10520,6415,13907,2189,4651,14040,9513,1348,14220,4561,302,502,15176,6243,4631,15617,14913,2419,12818,11676,118,12038,15688,1000,10822,9955,1763,5050,10118,8788,5812,6684,9210,13600,3582,1743,10284,3319,7836,1590,1947,6077,3578,14314,4814,12076},
			expectedBlockProposer:        1888,
		},
		{
			name:                         "slot 36, pool id 2, committee id 2",
			epoch:                        1,
			slot:						  36,
			poolId:                       2,
			committeeId: 				  2,
			expectedVaultCommittee:       []uint64{6609,3774,10541,13581,14271,11234,2606,3129,1411,2310,14341,6970,13797,10264,3754,1606,8801,15894,2963,7971,3011,2560,2629,12973},
			expectedAttestationCommittee: []uint64{12690,10380,14995,15746,4544,15942,5545,7975,6888,10062,10113,7555,13862,11403,3426,11398,183,2191,9942,4908,1463,9774,14485,1361,7349,5533,14469,828,5249,4861,9753,13429,15007,361,10242,6692,14562,4475,5097,9524,1317,8126,9957,7909,12347,9699,2420,4520,14146,4641,5238,8550,14211,10932,7176,13029,9152,15162,4883,2331,475,3989,30,15328,2416,15021,3765,7389,11351,12743,13450,9638,5135,11328,3381,10502,14150,11502,12002,2692,10066,8869,3894,228,1685,14561,9728,4097,2830,15835,3976,13548,9151,6287,1820,2807,7193,3675,15681,5733,11568,16103,7388,15522,9696,7471,7398,12133,6230,12277,2174,9623,11939,8474,12891,1393,5925,8897,834,5421,5636,12569,13944,8308,9602,4400,2615,12023},
			expectedBlockProposer:        16024,
		},
		{
			name:                         "slot 37, pool id 3, committee id 3",
			epoch:                        2,
			slot:						  37,
			poolId:                       3,
			committeeId: 				  3,
			expectedVaultCommittee:       []uint64{13980,9807,7798,15850,4252,15298,13026,12785,9445,9847,7417,11279,7269,14368,12496,13794,3068,7238,11969,13285,9872,4933,1731,8684},
			expectedAttestationCommittee: []uint64{14781,14908,8104,2778,6931,5126,1240,8317,15075,4706,7337,11099,13341,12795,11725,11343,9118,11674,5612,10054,4453,11141,10663,11181,1554,3086,11669,14362,1668,8162,3109,7845,6627,6270,1773,3363,3247,1436,2024,9791,4710,14179,9773,12397,8340,2576,4650,15894,2224,9251,11880,6222,10065,14658,1973,11152,2339,2235,7511,2847,9178,8031,709,5115,309,12840,13092,9200,742,10647,7705,3048,9787,11606,5920,5007,13000,4233,13916,7253,9705,11218,2804,3273,6256,1808,11436,9100,85,9262,7037,6862,2895,8977,12401,8114,10045,4743,12095,8844,12473,4899,3336,4386,13388,13988,9459,1751,14140,5527,2815,4851,3361,12865,12875,10549,16308,8892,6205,11382,12546,14765,10921,2232,3685,6842,13209,15858,},
			expectedBlockProposer:        15312,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pc,err := VaultCommittee(state, test.poolId, test.epoch)
			require.NoError(t, err)
			require.EqualValues(t, test.expectedVaultCommittee, pc)

			voting,err := SlotCommittee(state, test.slot, test.committeeId)
			require.NoError(t, err)
			require.EqualValues(t, test.expectedAttestationCommittee, voting)

			proposer,err := BlockProposer(state, test.slot)
			require.NoError(t, err)
			require.EqualValues(t, test.expectedBlockProposer, proposer)
		})
	}
}
