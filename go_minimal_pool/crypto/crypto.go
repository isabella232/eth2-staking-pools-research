package crypto

//func Random32Bytes() ([32]byte, error) {
//	s := [32]byte{}
//	_, err := rand.Read(s[:])
//	return s,err
//}
//
//func IntPow(a uint32, b uint32) uint32 {
//	return uint32(math.Pow(float64(a),float64(b)))
//}
//
//func IntToByteArray(num uint32) []byte {
//	size := int(unsafe.Sizeof(num))
//	arr := make([]byte, size)
//	for i := 0 ; i < size ; i++ {
//		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
//		arr[i] = byt
//	}
//	return arr
//}
//
//func ByteArrayToInt(arr []byte) uint32{
//	val := uint32(0)
//	size := len(arr)
//	for i := 0 ; i < size ; i++ {
//		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
//	}
//	return val
//}