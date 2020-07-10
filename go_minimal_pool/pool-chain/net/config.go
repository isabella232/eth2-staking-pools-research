package net

import (
	"encoding/hex"
	"time"
)

type NetworkConfig struct {
	PoolSize uint8
	NumberOfPools uint8
	PoolThreshold uint8

	SeedShuffleRoudnCount uint8

	EpochSpanSec time.Duration
	EpochTestMessage []byte

	GenesisSeed [32]byte // used for random beacon
}

func NewTestNetworkConfig() *NetworkConfig {
	_seed, _ := hex.DecodeString("b581262ce281d1e9deaf2f0158d7cd05217f1196d95956c5f55d837ccc3c8a9")
	var seed [32]byte
	copy(seed[:], _seed)


	_testMsg, _ := hex.DecodeString("292ea14188b703cbd3efde48f0952b15c4cc6c254f221e5669709888ccfbf8bf") // sha256 of 'test epoch msg'

	return &NetworkConfig{
		PoolSize:      3,
		PoolThreshold: 3,
		NumberOfPools: 2,
		SeedShuffleRoudnCount: 10,
		EpochSpanSec:  time.Second * 8,
		EpochTestMessage: _testMsg,
		GenesisSeed:   seed,
	}
}

func (c *NetworkConfig) TotalNumberOfParticipants() uint32 {
	return uint32(c.NumberOfPools * c.PoolSize)
}

func (c *NetworkConfig) ParticipantIndexesList() []uint32 {
	s := make([]uint32, c.TotalNumberOfParticipants())
	start := uint32(1)
	for i := range s {
		s[i] = start
		start += 1
	}
	return s
}
