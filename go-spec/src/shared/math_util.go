package shared

// Max returns the larger integer of the two
// given ones.This is used over the Max function
// in the standard math library because that max function
// has to check for some special floating point cases
// making it slower by a magnitude of 10.
func Max(a uint64, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

// Min returns the smaller integer of the two
// given ones. This is used over the Min function
// in the standard math library because that min function
// has to check for some special floating point cases
// making it slower by a magnitude of 10.
func Min(a uint64, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}