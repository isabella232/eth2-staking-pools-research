package src

type BeaconDuty struct {
	Type				uint8 // 0 - attestation, 1 - block proposal
	Slot				uint64
	Included			bool // whether or not it was included in the beacon chain (the pool earned reward from it)
	Executors			[16]byte // 128 bit of the executors (by order) which executed this duty
}


