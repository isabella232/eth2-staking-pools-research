package block

type BeaconDuty struct {
	DutyType      uint8 // 0 - attestation, 1 - block proposal
	Committee     uint64
	Slot          uint64
	Finalized     bool     // whether or not it was included in the beacon chain (the pool earned reward from it)
	Participation [16]byte // 128 bit of the executors (by order) which executed this duty
}

func NewBeaconDuty(
	dutyType      uint8,
	committee     uint64,
	slot          uint64,
	finalized     bool,
	participation [16]byte,
	) *BeaconDuty {
	return &BeaconDuty{
		DutyType:      dutyType,
		Committee:     committee,
		Slot:          slot,
		Finalized:     finalized,
		Participation: participation,
	}
}

func (duty *BeaconDuty) GetType() uint8 {
	return duty.DutyType
}

func (duty *BeaconDuty) GetCommittee() uint64 {
	return duty.Committee
}

func (duty *BeaconDuty) GetSlot() uint64 {
	return duty.Slot
}

func (duty *BeaconDuty) IsFinalized() bool  {
	return duty.Finalized
}

func (duty *BeaconDuty) GetParticipation() [16]byte  {
	return duty.Participation
}
