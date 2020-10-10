package params

import (
	"encoding/hex"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
)

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
	}
}

func UseMinimalTestConfig() {
	ChainConfig = testConfig()
}

// utils func
func SlotToEpoch(slot uint64) uint64 {
	return slot/ ChainConfig.SlotsInEpoch
}