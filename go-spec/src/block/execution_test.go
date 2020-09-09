package block

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/mocks"
	"github.com/golang/mock/gomock"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func MockedState(t *testing.T, participatedIndexes map[uint64]bool) (*mocks.MockIState, *gomock.Controller) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	ctrl := gomock.NewController(t)

	state := mocks.NewMockIState(ctrl)

	//
	pools := make([]core.IPool, 1)

	// BPs
	bps := make([]core.IBlockProducer, len(pools) * int(core.TestConfig().PoolExecutorsNumber))
	for i := 0 ; i < len(bps) ; i++ {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		bp := mocks.NewMockIBlockProducer(ctrl)
		bps[i] = bp

		bp.EXPECT().GetId().Return(uint64(i))
		if participatedIndexes[uint64(i)] {
			bp.EXPECT().IncreaseBalance(uint64(100))
		} else {
			bp.EXPECT().DecreaseBalance(uint64(100))
		}

		state.EXPECT().GetBlockProducer(uint64(i)).Return(bp)
	}

	// pools
	for i := 0 ; i < len(pools) ; i++ {
		executors := make([]uint64, core.TestConfig().PoolExecutorsNumber)
		for j := 0; j < int(core.TestConfig().PoolExecutorsNumber); j++ {
			executors[j] = bps[i*int(core.TestConfig().PoolExecutorsNumber)+j].GetId()
		} // no need to sort as they are already

		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		pool := mocks.NewMockIPool(ctrl)
		pools[i] = pool

		pool.EXPECT().GetSortedExecutors().Return(executors)
		state.EXPECT().GetPool(uint64(i)).Return(pool)
	}


	return state, ctrl
}

func TestAttestationSuccessful(t *testing.T) {
	// create mocked state
	state, ctrl := MockedState(t, map[uint64]bool{0:true})
	defer ctrl.Finish()

	// generate summary
	summary := NewExecutionSummary(
			0,
			1,
			[]*BeaconDuty{
				NewBeaconDuty(
						0,
						0,
						0,
						true,
						[16]byte{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}, // the first executor index is set to 1
					),
			},
		)

	require.NoError(t, summary.ApplyOnState(state))

	//for _, duty := range summary.GetDuties() {
//		pool := state.GetPool(summary.PoolId)
//
//		for i:=0 ; i < int(core.TestConfig().PoolExecutorsNumber) ; i++ {
//			bp := state.GetBlockProducer(pool.GetSortedExecutors()[i])
//
//			Participation := duty.GetParticipation()
//			if shared.IsBitSet(Participation[:], uint64(i)) {
//				require.EqualValues(t, 1100, bp.GetBalance())
//				require.EqualValues(t, 0, i)
//			} else {
//				require.EqualValues(t, 900, bp.GetBalance())
//			}
//		}
//	}
}