package pool_chain

import (
	"encoding/hex"
	"time"
)

type NetworkConfig struct {
	PoolSize uint8
	NumberOfPools uint8

	EpochSpanSec time.Duration

	GenesisSeed []byte // used for random beacon
}

func NewTestNetworkConfig() *NetworkConfig {
	seed, _ := hex.DecodeString("")

	return &NetworkConfig{
		PoolSize:      3,
		NumberOfPools: 1,
		EpochSpanSec:  time.Second * 8,
		GenesisSeed:   seed,
	}
}

func (c *NetworkConfig) TotalNumberOfParticipants() uint32 {
	return uint32(c.NumberOfPools * c.PoolSize)
}
