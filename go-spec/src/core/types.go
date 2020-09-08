package core

import "github.com/herumi/bls-eth-go-binary/bls"

type IState interface {
	Copy() (IState, error)
	Root() ([32]byte,error)
	GetPools() []IPool
	GetPool(id uint64) IPool
	GetBlockProducers() []IBlockProducer
	GetBlockProducer(id uint64) IBlockProducer
	GetCurrentEpoch() uint64
	GetSeed() [32]byte
	GetPastSeed(epoch uint64) [32]byte
	AddNewPool(pool IPool) error

	// For a given epoch and pool id, return the pool's committee
	PoolExecutors(poolId uint64, epoch uint64) ([]uint64, error)
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
	GetId() uint64
	GetPubKey() *bls.PublicKey
	GetSortedExecutors() [16]uint64
}

type IBlockBody interface {
	GetEpochNumber() uint64
	GetProposer() uint64
	GetExecutionSummaries() []IExecutionSummary
	GetNewPoolRequests() []ICreatePoolRequest
	GetStateRoot() []byte
	GetParentBlockRoot() []byte
	Root() []byte
	Validate() error
}

type IBlockHeader interface {
	Copy() IBlockHeader
	Validate(bp IBlockProducer) error
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
	GetBalance() uint64
	GetStake() uint64
	IsSlashed() bool
	IsActive() bool
	IncreaseBalance(id uint64, change uint64) (newBalance uint64, error error)
	DecreaseBalance(id uint64, change uint64) (newBalance uint64, error error)
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
