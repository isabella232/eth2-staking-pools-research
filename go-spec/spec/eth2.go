package spec

type BeaconDuty struct {
	Type				uint8 // 0 - attestation, 1 - block proposal
	Reward				uint64
}

type Performance struct {
	Execution map[BeaconDuty][16]byte // for every duty specify an array of 128 bits (16 bytes) of who participated in the execution of that duty
}


