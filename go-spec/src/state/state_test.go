package state

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/block"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func GenerateAttestationSuccessfulSummary() testing.InternalExample {
	return &PoolExecutionSummary{
		PoolId:        0,
		Epoch:         1,
		Duties:        []*BeaconDuty{
			&BeaconDuty{
				dutyType:      0,
				slot:          0,
				finalized:     true,
				participation: [16]byte{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}, // the first executor index is set to 1
			},
		},
	}
}

func GenerateRandomState(t *testing.T) *State {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	pools := make([]core.IPool, 5)

	//
	bps := make([]core.IBlockProducer, len(pools) * int(core.TestConfig().PoolExecutorsNumber))
	for i := 0 ; i < len(bps) ; i++ {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		bps[i] = &BlockProducer{
			id:      uint64(i),
			balance: 1000,
			stake:   0,
			slashed: false,
			active:  true,
			pubKey:  sk.GetPublicKey(),
		}
	}

	//
	for i := 0 ; i < len(pools) ; i++ {
		executors := make([]uint64, core.TestConfig().PoolExecutorsNumber)
		for j := 0 ; j < int(core.TestConfig().PoolExecutorsNumber) ; j++ {
			executors[j] = bps[i*int(core.TestConfig().PoolExecutorsNumber) + j].GetId()
		} // no need to sort as they are already

		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		pools[i] = &Pool{
			id:              uint64(i),
			sortedExecutors: executors,
			pubKey:          sk.GetPublicKey(),
		}
	}

	return &State{
		pools: pools,
		headBlockHeader: &block.BlockHeader{
			BlockRoot: nil,
			Signature: nil,
		},
		blockProducers: bps,
		seed:           shared.SliceToByte32([]byte("seedseedseedseedseedseedseedseed")),
	}
}

func TestRandaoSeedMix(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	sk := &bls.SecretKey{}
	sk.SetByCSPRNG()

	state := GenerateRandomState(t)

	body, err := block.NewBlockBody(0, 1, state, make([]core.IExecutionSummary, 0), []byte("parent"))
	require.NoError(t, err)
	header,err := block.NewBlockHeader(sk, body)
	require.NoError(t, err)

	newState, err := state.ProcessNewBlock(header, body)
	require.NoError(t, err)

	expectedSeed,err := shared.MixSeed(state.seed, shared.SliceToByte32(header.Signature))
	require.NoError(t, err)
	require.EqualValues(t, expectedSeed, newState.GetSeed())
}


func TestBlockValidation(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	sk := &bls.SecretKey{}
	sk.SetByCSPRNG()

	// generate header and body
	block.NewBlockBody(
		0,
		[]core.IExecutionSummary{generate}
		)
	body := &block.BlockBody{
		proposer:           0,
		executionSummaries: []*block.PoolExecutionSummary{src.GenerateAttestationSuccessfulSummary()},
		newPoolReq:			   []*block.CreatePoolRequest{
			&block.CreatePoolRequest{
				Id:                  0,
				Status:              0,
				StartEpoch:          0,
				EndEpoch:            1,
				LeaderBlockProducer: 0,
				CreatedPubKey:       nil,
				Participation:       [16]byte{},
			},
		},
		stateRoot:       []byte("root"),
		parentBlockRoot: []byte("parent"),
	}
	root,err := body.Root()
	require.NoError(t, err)

	// sign
	sig := sk.SignByte(root[:])

	head := &block.BlockHeader{
		BlockRoot: root[:],
		Signature: sig.Serialize(),
	}

	state := GenerateRandomState(t)

	// set BP
	bp, err := state.GetBlockProducer(0)
	require.NoError(t, err)
	bp.PubKey = sk.GetPublicKey()

	require.NoError(t, state.ValidateBlock(head, body))
}

func TestCreatedNewPoolReq(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	sk := &bls.SecretKey{}
	sk.SetByCSPRNG()

	reqs := []*block.CreatePoolRequest{
		&block.CreatePoolRequest{
			Id:                  0,
			Status:              1,
			StartEpoch:          0,
			EndEpoch:            1,
			LeaderBlockProducer: 0,
			CreatedPubKey:       sk.GetPublicKey().Serialize(),
			Participation:       [16]byte{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}, // pos 0 is set only
		},
	}

	helperFunc = src.NewSimpleFunctions()
	state := GenerateRandomState(t)
	currentBP, err := state.GetBlockProducer(0)
	require.NoError(t, err)

	// save current state for fetching
	require.NoError(t, helperFunc.SaveState(state, 0))

	participants,err := state.DKGParticipants(0)

	require.NoError(t, state.ProcessNewPoolRequests(reqs, currentBP))
	require.Equal(t, 6, len(state.pools))

	// check new balances
	currentBP, err = state.GetBlockProducer(currentBP.Id)
	require.NoError(t, err)
	require.EqualValues(t, 4000, currentBP.Balance)

	bp, err := state.GetBlockProducer(participants[0])
	require.NoError(t, err)
	require.EqualValues(t, 2000, bp.Balance)

	for i := 1 ; i < len(participants) ; i++ {
		bp, err := state.GetBlockProducer(participants[i])
		require.NoError(t, err)
		require.EqualValues(t, 0, bp.Balance)
	}
}

func TestFailedToCreateNewPool(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	sk := &bls.SecretKey{}
	sk.SetByCSPRNG()

	reqs := []*block.CreatePoolRequest{
		&block.CreatePoolRequest{
			Id:                  0,
			Status:              2,
			StartEpoch:          0,
			EndEpoch:            1,
			LeaderBlockProducer: 0,
			CreatedPubKey:       sk.GetPublicKey().Serialize(),
			Participation:       [16]byte{1,1,1,1,0,0,0,0,0,0,1,0,0,0,0,0}, // pos 0 is set only
		},
	}

	helperFunc = src.NewSimpleFunctions()
	state := GenerateRandomState(t)
	currentBP, err := state.GetBlockProducer(0)
	require.NoError(t, err)

	// save current state for fetching
	require.NoError(t, helperFunc.SaveState(state, 0))

	participants,err := state.DKGParticipants(0)

	require.NoError(t, state.ProcessNewPoolRequests(reqs, currentBP))
	require.Equal(t, 5, len(state.pools))

	// check new balances
	currentBP, err = state.GetBlockProducer(currentBP.Id)
	require.NoError(t, err)
	require.EqualValues(t, 1000, currentBP.Balance)

	for i := 0 ; i < len(participants) ; i++ {
		bp, err := state.GetBlockProducer(participants[i])
		require.NoError(t, err)
		require.EqualValues(t, 0, bp.Balance)
	}
}