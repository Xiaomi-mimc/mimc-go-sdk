package byteutil

import (
	"bytes"
	"hash/adler32"
)

func ToBytes(int8s []int8) []byte {
	size := len(int8s)
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bytes[i] = byte(int8s[i])
	}
	return bytes
}

func Copy(src *[]byte, from, length int) []byte {
	dst := make([]byte, length)
	for i := 0; i < length; i++ {
		dst[i] = (*src)[from+i]
	}
	return dst
}

func Bytes(data int) []byte {
	bytes := make([]byte, 4)
	first := byte((data >> 24) & 0xff)
	second := byte((data >> 16) & 0xff)
	third := byte((data >> 8) & 0xff)
	fourth := byte(data & 0xff)
	bytes[0] = first
	bytes[1] = second
	bytes[2] = third
	bytes[3] = fourth
	return bytes
}

func TransferUint16(bytes *[]byte, data uint16, index int) {
	low := byte(data & 0xff)
	high := byte((data >> 8) & 0xff)
	(*bytes)[index] = high
	(*bytes)[index+1] = low
}
func GetUint16FromBytes(bytes *[]byte, index int) uint16 {
	firstByte, secondByte := 0, 0
	firstByte = 0x000000FF & int((*bytes)[index])
	secondByte = 0x000000FF & int((*bytes)[index+1])
	return uint16((firstByte<<8 | secondByte) & 0xFFFFFFFF)
}
func GetIntFromBytes(bytes *[]byte, index int) int {
	firstByte, secondByte, thirdByte, fourthByte := 0, 0, 0, 0
	firstByte = 0x000000FF & int((*bytes)[index])
	secondByte = 0x000000FF & int((*bytes)[index+1])
	thirdByte = 0x000000FF & int((*bytes)[index+2])
	fourthByte = 0x000000FF & int((*bytes)[index+3])
	return ((firstByte<<24 | secondByte<<16 | thirdByte<<8 | fourthByte) & 0xFFFFFFFF)
}

func TransferInt(bytes *[]byte, data int, index int) {
	first := byte((data >> 24) & 0xff)
	second := byte((data >> 16) & 0xff)
	third := byte((data >> 8) & 0xff)
	fourth := byte(data & 0xff)
	(*bytes)[index] = first
	(*bytes)[index+1] = second
	(*bytes)[index+2] = third
	(*bytes)[index+3] = fourth
}

func Crc(bytes []byte) int {
	return int(adler32.Checksum(bytes))
}

func Integrate(byteArray ...[]byte) []byte {
	size := len(byteArray)
	buffer := new(bytes.Buffer)
	for i := 0; i < size; i++ {
		buffer.Write(byteArray[i])
	}
	return buffer.Bytes()
}
