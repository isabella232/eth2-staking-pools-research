package core

import "github.com/herumi/bls-eth-go-binary/bls"

type IState interface {
	Copy() (IState, error)
	Root() ([32]byte,error)
	GetPools() []IPool
	GetPool(id uint64) IPool
	AddNewPool(pool IPool) error
	GetBlockProducers() []IBlockProducer
	GetBlockProducer(id uint64) IBlockProducer
	GetCurrentEpoch() uint64
	SetCurrentEpoch(epoch uint64)
	GetHeadBlockHeader() IBlockHeader
	SetHeadBlockHeader(header IBlockHeader)
	GetSeed(epoch uint64) [32]byte
	SetSeed(seed [32]byte, epoch uint64)
	GetPastSeed(epoch uint64) [32]byte
	GetBlockRoot(epoch uint64) [32]byte
	SetBlockRoot(root [32]byte, epoch uint64)
	GetStateRoot(epoch uint64) [32]byte
	SetStateRoot(root [32]byte, epoch uint64)

	// For a given epoch and pool id, return the pool's committee
	PoolCommittee(poolId uint64, epoch uint64) ([]uint64, error)
	// For a given epoch and create pool request id, return the DKG committee
	DKGCommittee(reqId uint64, epoch uint64)([]uint64, error)
	// For a given epoch, return the voting committee
	BlockVotingCommittee(epoch uint64)([]uint64, error)
	// For a given epoch, return the block proposer
	GetBlockProposer(epoch uint64) (uint64, error)

	ValidateBlock(header IBlockHeader, body IBlockBody) error
	ProcessPoolExecutions(summaries []IExecutionSummary) error
	ProcessNewPoolRequests(requests []ICreatePoolRequest) error
	ProcessNewBlock(newBlockHeader IBlockHeader, newBlockBody IBlockBody) (newState IState, error error)
}

type IPool interface {
	Copy() (IPool, error)
	IsActive() bool
	SetActive(status bool)
	GetId() uint64
	GetPubKey() *bls.PublicKey
	GetSortedExecutors() []uint64
	SetSortedExecutors(executors []uint64)
}

type IBlockBody interface {
	GetEpochNumber() uint64
	GetProposer() uint64
	GetExecutionSummaries() []IExecutionSummary
	GetNewPoolRequests() []ICreatePoolRequest
	GetStateRoot() []byte
	GetParentBlockRoot() []byte
	Root() ([]byte, error)
	Validate() error
}

type IBlockHeader interface {
	Copy() IBlockHeader
	Validate(bp IBlockProducer) error
	GetBlockRoot() []byte
	GetSignature() []byte
}

type IBeaconDuty interface {
	GetType() uint8 // 0 - attestation, 1 - proposal, 2 - aggregation
	GetCommittee() uint64
	GetSlot() uint64
	IsFinalized() bool // marked true if that duty was finalized on the beacon chain
	GetParticipation() [16]byte // 128 bit of the executors (by order) which executed this duty
}

type IExecutionSummary interface {
	GetPoolId() uint64
	GetEpoch() uint64
	GetDuties() []IBeaconDuty
	ApplyOnState(state IState) error
}

type IBlockProducer interface {
	Copy() (IBlockProducer, error)
	GetId() uint64
	GetPubKey() *bls.PublicKey
	SetPubKey(pk *bls.PublicKey)
	GetBalance() uint64
	GetStake() uint64
	IsSlashed() bool
	IsActive() bool
	SetExited(atEpoch uint64) // will mark active=false ad exit epoch = atEpoch
	ExitEpoch() uint64 // will return 0 if is still active
	IncreaseBalance(change uint64) (newBalance uint64, error error)
	DecreaseBalance(change uint64) (newBalance uint64, error error)
}

type ICreatePoolRequest interface {
	GetId() uint64
	GetStatus() uint64
	GetStartEpoch() uint64
	GetEndEpoch() uint64
	GetLeaderBP() uint64
	GetCreatePubKey() []byte
	GetParticipation() [16]byte // 128 bit of the committee (by order) which executed this request
	Validate(state IState, currentBP IBlockProducer) error
}
