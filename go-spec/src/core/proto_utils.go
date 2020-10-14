package core

import "bytes"

func CheckpointsEqual(l *Checkpoint, r *Checkpoint) bool {
	return l.Epoch == r.Epoch && bytes.Equal(l.Root, r.Root)
}

// returns true if equal
func AttestationDataEqual(att1 *AttestationData, att2 *AttestationData) bool {
	return att1.Slot == att2.Slot &&
		CheckpointsEqual(att1.Target, att2.Target) &&
		CheckpointsEqual(att1.Source, att2.Source) &&
		bytes.Equal(att1.BeaconBlockRoot, att2.BeaconBlockRoot)
}