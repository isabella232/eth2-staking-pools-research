package shared

func SliceToByte32(slice []byte) [32]byte {
	var arr [32]byte
	copy(arr[:], slice[:32])
	return arr
}
