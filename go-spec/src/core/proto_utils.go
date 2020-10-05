package core

import "bytes"

func CheckpointsEqual(l *Checkpoint, r *Checkpoint) bool {
	return l.Epoch == r.Epoch && bytes.Equal(l.Root, r.Root)
}
