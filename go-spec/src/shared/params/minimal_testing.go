package params

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
)

// Bytes4 returns integer x to bytes in little-endian format, x.to_bytes(4, 'little').
// TODO - copied here for cyclic dependency issue
func Bytes4(x uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, x)
	return bytes[:4]
}

func testConfig() *core.PoolsChainConfig {
	genesisSeed,_ := hex.DecodeString("sdddseedseedseedseedseedseedseed")

	return &core.PoolsChainConfig{
		GenesisSeed: 	       genesisSeed,
		GenesisEpoch: 		   0,

		VaultSize:             4,
		BaseEth2DutyReward:    100,
		DKGReward:             1000,

		SlotsInEpoch:                32,
		MinAttestationCommitteeSize: 16,
		MaxAttestationCommitteeSize: 16,
		MinAttestationInclusionDelay: 1,

		ZeroHash: make([]byte, 32),

		DomainBeaconProposer: Bytes4(0),
		DomainBeaconAttester: Bytes4(1),
		DomainRandao: Bytes4(2),
		GenesisForkVersion: []byte{},
	}
}

func UseMinimalTestConfig() {
	ChainConfig = testConfig()
}

// utils func
func SlotToEpoch(slot uint64) uint64 {
	return slot/ ChainConfig.SlotsInEpoch
}