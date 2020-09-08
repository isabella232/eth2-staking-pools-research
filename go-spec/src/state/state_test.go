package state

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func GenerateRandomState(t *testing.T) *State {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	pools := make([]*Pool, 5)

	//
	bps := make([]*BlockProducer, len(pools) * int(core.TestConfig().PoolExecutorsNumber))
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
			executors[j] = bps[i*int(core.TestConfig().PoolExecutorsNumber) + j].id
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
		HeadBlockHeader: &src.BlockHeader{
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

	body := &src.BlockBody{
		Proposer:              0,
		Number:                1,
		PoolsExecutionSummary: make([]*src.PoolExecutionSummary, 0),
		StateRoot:             []byte("state"),
		ParentBlockRoot:       []byte("parent"),
	}
	header,err := src.NewBlockHeader(sk, body)
	require.NoError(t, err)

	helperFunc = src.NewSimpleFunctions()
	newState, err := state.ProcessNewBlock(header, body)
	require.NoError(t, err)

	expectedSeed,err := shared.MixSeed(state.seed, shared.SliceToByte32(header.Signature))
	require.NoError(t, err)
	require.EqualValues(t, expectedSeed, newState.seed)
}


func TestBlockValidation(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	sk := &bls.SecretKey{}
	sk.SetByCSPRNG()

	// generate header and body
	body := &src.BlockBody{
		Proposer:              0,
		PoolsExecutionSummary: []*src.PoolExecutionSummary{src.GenerateAttestationSuccessfulSummary()},
		NewPoolReq:			   []*src.CreatePoolRequest{
			&src.CreatePoolRequest{
				Id:                  0,
				Status:              0,
				StartEpoch:          0,
				EndEpoch:            1,
				LeaderBlockProducer: 0,
				CreatedPubKey:       nil,
				Participation:       [16]byte{},
			},
		},
		StateRoot:             []byte("root"),
		ParentBlockRoot:       []byte("parent"),
	}
	root,err := body.Root()
	require.NoError(t, err)

	// sign
	sig := sk.SignByte(root[:])

	head := &src.BlockHeader{
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

	reqs := []*src.CreatePoolRequest{
		&src.CreatePoolRequest{
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

	reqs := []*src.CreatePoolRequest{
		&src.CreatePoolRequest{
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