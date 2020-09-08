package block

type BeaconDuty struct {
	dutyType      uint8 // 0 - attestation, 1 - block proposal
	committee     uint64
	slot          uint64
	finalized     bool     // whether or not it was included in the beacon chain (the pool earned reward from it)
	participation [16]byte // 128 bit of the executors (by order) which executed this duty
}

func (duty *BeaconDuty) GetType() uint8 {
	return duty.dutyType
}

func (duty *BeaconDuty) GetCommittee() uint64 {
	return duty.committee
}

func (duty *BeaconDuty) GetSlot() uint64 {
	return duty.slot
}

func (duty *BeaconDuty) IsFinalized() bool  {
	return duty.finalized
}

func (duty *BeaconDuty) GetParticipation() [16]byte  {
	return duty.participation
}
