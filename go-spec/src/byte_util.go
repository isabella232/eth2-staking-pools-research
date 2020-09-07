package src

import "encoding/binary"

// takes a src byte array and sets the bit at position pos to 1
func SetBit(src []byte, pos uint64, size uint64) []byte {
	mask := uint64(1 << pos)
	srcNum := binary.LittleEndian.Uint64(src)

	ret := make([]byte, size)
	binary.LittleEndian.PutUint64(ret, srcNum | mask)
	return ret
}

func IsBitSet (src []byte, pos uint64) bool {
	mask := uint64(1 << pos)
	srcNum := binary.LittleEndian.Uint64(src)
	return (srcNum & mask) != 0
}

func SliceToByte32(slice []byte) [32]byte {
	var arr [32]byte
	copy(arr[:], slice[:32])
	return arr
}