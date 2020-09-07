package src

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func GenerateRandomState(t *testing.T) *State {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	pools := make([]*Pool, 5)

	//
	bps := make([]*BlockProducer, len(pools) * int(TestConfig().PoolExecutorsNumber))
	for i := 0 ; i < len(bps) ; i++ {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		bps[i] = &BlockProducer{
			Id:      uint64(i),
			Balance: 1000,
			Stake:   0,
			Slashed: false,
			Active: true,
			PubKey: sk.GetPublicKey(),
		}
	}

	//
	for i := 0 ; i < len(pools) ; i++ {
		executors := make([]uint64, TestConfig().PoolExecutorsNumber)
		for j := 0 ; j < int(TestConfig().PoolExecutorsNumber) ; j++ {
			executors[j] = bps[i*int(TestConfig().PoolExecutorsNumber) + j].Id
		} // no need to sort as they are already

		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		pools[i] = &Pool{
			Id:              uint64(i),
			SortedExecutors: executors,
			PubKey: sk.GetPublicKey(),
		}
	}

	return &State{
		Pools:           pools,
		HeadBlockHeader: &BlockHeader{
			BlockRoot: nil,
			Signature: nil,
		},
		BlockProducers:  bps,
		Seed:            SliceToByte32([]byte("seedseedseedseedseedseedseedseed")),
	}
}

func TestRandaoSeedMix(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	sk := &bls.SecretKey{}
	sk.SetByCSPRNG()

	state := GenerateRandomState(t)

	body := &BlockBody{
		Proposer:              0,
		Number:                1,
		PoolsExecutionSummary: make([]*PoolExecutionSummary, 0),
		StateRoot:             []byte("state"),
		ParentBlockRoot:       []byte("parent"),
	}
	header,err := NewBlockHeader(sk, body)
	require.NoError(t, err)

	helperFunc = NewSimpleFunctions()
	newState, err := state.ProcessNewBlock(header, body)
	require.NoError(t, err)

	expectedSeed,err := MixSeed(state.Seed, SliceToByte32(header.Signature))
	require.NoError(t, err)
	require.EqualValues(t, expectedSeed, newState.Seed)
}


func TestBlockValidation(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	sk := &bls.SecretKey{}
	sk.SetByCSPRNG()

	// generate header and body
	body := &BlockBody{
		Proposer:              0,
		PoolsExecutionSummary: []*PoolExecutionSummary{GenerateAttestationSuccessfulSummary()},
		NewPoolReq:			   []*CreatePoolRequest{
			&CreatePoolRequest{
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

	head := &BlockHeader{
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

	reqs := []*CreatePoolRequest{
		&CreatePoolRequest{
			Id:                  0,
			Status:              1,
			StartEpoch:          0,
			EndEpoch:            1,
			LeaderBlockProducer: 0,
			CreatedPubKey:       sk.GetPublicKey().Serialize(),
			Participation:       [16]byte{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}, // pos 0 is set only
		},
	}

	helperFunc = NewSimpleFunctions()
	state := GenerateRandomState(t)
	currentBP, err := state.GetBlockProducer(0)
	require.NoError(t, err)

	// save current state for fetching
	require.NoError(t, helperFunc.SaveState(state, 0))

	participants,err := state.DKGParticipants(0)

	require.NoError(t, state.ProcessNewPoolRequests(reqs, currentBP))
	require.Equal(t, 6, len(state.Pools))

	// check new balances
	bp, err := state.GetBlockProducer(participants[0])
	require.NoError(t, err)
	require.EqualValues(t, 2000, bp.Balance)

	for i := 1 ; i < len(participants) ; i++ {
		bp, err := state.GetBlockProducer(participants[i])
		require.NoError(t, err)
		require.EqualValues(t, 0, bp.Balance)
	}
}