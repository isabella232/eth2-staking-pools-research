package params

import (
	"encoding/hex"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
)

func testConfig() *core.PoolsChainConfig {
	genesisSeed,_ := hex.DecodeString("sdddseedseedseedseedseedseedseed")

	return &core.PoolsChainConfig{
		GenesisSeed: 	       genesisSeed,

		VaultSize:             24,
		BaseEth2DutyReward:    100,
		DKGReward:             1000,

		SlotsInEpoch:                32,
		MinAttestationCommitteeSize: 128,
		MaxAttestationCommitteeSize: 2048,
		MinAttestationInclusionDelay: 1,
	}
}

func UseMinimalTestConfig() {
	ChainConfig = testConfig()
}

// utils func
func SlotToEpoch(slot uint64) uint64 {
	return slot/ ChainConfig.SlotsInEpoch
}