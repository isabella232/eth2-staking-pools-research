package spec

type BeaconDuty struct {
	Type				uint8 // 0 - attestation, 1 - block proposal
	Slot				uint64
	Included			bool // whether or not it was included in the beacon chain (the pool earned reward from it)
}


