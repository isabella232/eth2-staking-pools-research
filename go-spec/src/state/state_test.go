package state

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/block"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

type DummyBeaconDuty struct {
	t uint8
	committee uint64
	slot uint64
	finalized bool
	participation [16]byte
}
func (d *DummyBeaconDuty) GetType() uint8 { return d.t}
func (d *DummyBeaconDuty) GetCommittee() uint64 { return d.committee}
func (d *DummyBeaconDuty) GetSlot() uint64 { return d.slot}
func (d *DummyBeaconDuty) IsFinalized() bool { return d.finalized}
func (d *DummyBeaconDuty) GetParticipation() [16]byte { return d.participation}

type DummyExecSummary struct {
	poolId uint64
	epoch uint64
	duties []core.IBeaconDuty
}
func (s *DummyExecSummary) GetPoolId() uint64 { return s.poolId }
func (s *DummyExecSummary) GetEpoch() uint64 { return s.epoch}
func (s *DummyExecSummary) GetDuties() []core.IBeaconDuty { return s.duties}
func (s *DummyExecSummary) ApplyOnState(state core.IState) error { return nil }

func successfulAttestationSummary() core.IExecutionSummary {
	return &DummyExecSummary{
		poolId: 0,
		epoch:  1,
		duties: []core.IBeaconDuty{
				&DummyBeaconDuty{
					t:             0,
					committee:     0,
					slot:          0,
					finalized:     true,
					participation: [16]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // the first executor index is set to 1,
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

	body, err := block.NewBlockBody(0, 1, state, make([]core.IExecutionSummary, 0), nil, []byte("parent"))
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

	state := GenerateRandomState(t)

	sk := &bls.SecretKey{}
	sk.SetByCSPRNG()

	// generate header and body
	body,err := block.NewBlockBody(
			0,
			0,
			state,
			[]core.IExecutionSummary{successfulAttestationSummary()},
			[]core.ICreatePoolRequest{
				block.NewCreatePoolRequest(
					0,0,0,1,0,nil,[16]byte{},
				),
			},
			[]byte("parent"),
		)
	require.NoError(t, err)
	root,err := body.Root()
	require.NoError(t, err)

	// sign
	sig := sk.SignByte(root[:])

	head := &block.BlockHeader{
		BlockRoot: root[:],
		Signature: sig.Serialize(),
	}

	// set BP
	bp := state.GetBlockProducer(0)
	require.NoError(t, err)
	bp.SetPubKey(sk.GetPublicKey())

	require.NoError(t, state.ValidateBlock(head, body))
}

func TestCreatedNewPoolReq(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	sk := &bls.SecretKey{}
	sk.SetByCSPRNG()

	reqs := []core.ICreatePoolRequest{
		block.NewCreatePoolRequest(
			0,
			0,
			0,
			1,
			0,
			sk.GetPublicKey().Serialize(),
			[16]byte{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}, // pos 0 is set only,
		),
	}

	state := GenerateRandomState(t)
	currentBP := state.GetBlockProducer(0)

	// save current state for fetching
	require.NoError(t, helperFunc.SaveState(state, 0))

	participants,err := state.DKGCommittee(0,0)

	require.NoError(t, state.ProcessNewPoolRequests(reqs))
	require.Equal(t, 6, len(state.pools))

	// check new balances
	currentBP = state.GetBlockProducer(currentBP.GetId())
	require.NoError(t, err)
	require.EqualValues(t, 4000, currentBP.GetBalance())

	bp := state.GetBlockProducer(participants[0])
	require.EqualValues(t, 2000, bp.GetBalance())

	for i := 1 ; i < len(participants) ; i++ {
		bp := state.GetBlockProducer(participants[i])
		require.EqualValues(t, 0, bp.GetBalance())
	}
}

func TestFailedToCreateNewPool(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	sk := &bls.SecretKey{}
	sk.SetByCSPRNG()

	reqs := []core.ICreatePoolRequest {
		block.NewCreatePoolRequest(
			0,
			2,
			0,
			1,
			0,
			sk.GetPublicKey().Serialize(),
			[16]byte{1,1,1,1,0,0,0,0,0,0,1,0,0,0,0,0}, // random assignments
		),
	}

	state := GenerateRandomState(t)
	currentBP := state.GetBlockProducer(0)

	// save current state for fetching
	require.NoError(t, helperFunc.SaveState(state, 0))

	participants,err := state.DKGCommittee(0, 0)
	require.NoError(t, err)

	require.NoError(t, state.ProcessNewPoolRequests(reqs))
	require.Equal(t, 5, len(state.pools))

	// check new balances
	currentBP = state.GetBlockProducer(currentBP.GetId())
	require.EqualValues(t, 1000, currentBP.GetBalance())

	for i := 0 ; i < len(participants) ; i++ {
		bp := state.GetBlockProducer(participants[i])
		require.EqualValues(t, 0, bp.GetBalance())
	}
}