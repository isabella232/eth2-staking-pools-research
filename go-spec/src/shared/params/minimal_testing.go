package params

import "github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"

func testConfig() *core.PoolsChainConfig {
	return &core.PoolsChainConfig{
		VaultSize:             24,

		BaseEth2DutyReward:    100,
		DKGReward:             1000,

		SlotsInEpoch:                32,
		MinAttestationCommitteeSize: 128,
		MaxAttestationCommitteeSize: 2048,
	}
}

func UseMinimalTestConfig() {
	ChainConfig = testConfig()
}

// utils func
func SlotToEpoch(slot uint64) uint64 {
	return slot/ ChainConfig.SlotsInEpoch
}