package core

type Config struct {
	VaultSize             uint64
	DKGParticipantsNumber uint64
	BaseEth2DutyReward    uint64
	DKGReward             uint64

	// slots and epochs
	SlotsInEpoch                uint64
	MinAttestationCommitteeSize uint64
	MaxAttestationCommitteeSize	uint64
}

func TestConfig() *Config {
	return &Config{
		VaultSize:             24,
		DKGParticipantsNumber: 128, // TODO - remove and use VaultSize

		BaseEth2DutyReward:    100,
		DKGReward:             1000,

		SlotsInEpoch:                32,
		MinAttestationCommitteeSize: 128,
		MaxAttestationCommitteeSize: 2048,
	}
}


// utils func
func (c *Config) SlotToEpoch(slot uint64) uint64 {
	return slot/ c.SlotsInEpoch
}