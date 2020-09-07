package src

import (
	"encoding/binary"
	"github.com/stretchr/testify/require"
	"testing"
)

func numberToBytes(n uint64) []byte {
	ret := make([]byte, 64)
	binary.LittleEndian.PutUint64(ret, n)
	return ret
}

func TestIsBitSet(t *testing.T) {
	tests := []struct{
		testName string
		src []byte
		pos uint64
		expected bool
	}{
		{
			testName: "empty array, should return false",
			src: numberToBytes(0),
			pos: 1,
			expected: false,
		},
		{
			testName: "bit set, should return true",
			src: numberToBytes(2),
			pos: 1,
			expected: true,
		},
		{
			testName: "many bits set, should return true",
			src: numberToBytes(349525), // 01010101010101010101
			pos: 0,
			expected: true,
		},
		{
			testName: "many bits set, should return true",
			src: numberToBytes(349525), // 01010101010101010101
			pos: 2,
			expected: true,
		},
		{
			testName: "many bits set, should return false",
			src: numberToBytes(349525), // 01010101010101010101
			pos: 1,
			expected: false,
		},
		{
			testName: "many bits set, should return false",
			src: numberToBytes(349525), // 01010101010101010101
			pos: 3,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			require.Equal(t, test.expected, IsBitSet(test.src, test.pos))
		})
	}
}

func TestSetBit(t *testing.T) {
	tests := []struct{
		testName string
		src []byte
		pos uint64
		expected []byte
	}{
		{
			testName: "empty byte manipulation pos 1",
			src: numberToBytes(0),
			pos: 1,
			expected: numberToBytes(2),
		},
		{
			testName: "empty byte manipulation pos 2",
			src: numberToBytes(0),
			pos: 2,
			expected: numberToBytes(4),
		},
		{
			testName: "empty byte manipulation pos 12",
			src: numberToBytes(0),
			pos: 12,
			expected: numberToBytes(4096),
		},
		{
			testName: "non empty byte manipulation pos 2",
			src: numberToBytes(2),
			pos: 2,
			expected: numberToBytes(6),
		},
		{
			testName: "non empty byte manipulation pos 1",
			src: numberToBytes(349525), // 01010101010101010101
			pos: 1,
			expected: numberToBytes(349527), // 01010101010101010111
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			res := SetBit(test.src, test.pos, 64)
			require.Equal(t,test.expected, res)
		})
	}
}
