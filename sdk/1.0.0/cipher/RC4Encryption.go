package cipher

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mimc-go-sdk/util/string"
)

type RC4Cryption struct {
	keylength int
	S         []int8
	the_i     int
	the_j     int
	next_j    int
}

func create() *RC4Cryption {
	rc4 := new(RC4Cryption)
	rc4.S = make([]int8, 256)
	rc4.the_i, rc4.the_j, rc4.next_j, rc4.keylength = 0, 0, -666, 8
	return rc4
}

func (this RC4Cryption) init() {
	this.the_i, this.the_j = 0, 0
}

func (this RC4Cryption) ksa1(n int, key []byte, printstats bool) {
	keylength := len(key)
	for i := 0; i < 256; i++ {
		this.S[i] = int8(i)
	}
	this.the_j = 0
	for this.the_i = 0; this.the_i < n; this.the_i++ {
		this.the_j = (this.the_j + int(this.posify(this.S[this.the_i])+this.posify(int8(key[this.the_i%keylength])))) % 256
		this.sswap(&(this.S), this.the_i, this.the_j)
	}
	if n != 256 {
		this.next_j = (this.the_j + int(this.posify(this.S[n])+this.posify(int8(key[n%keylength])))) % 256
	}
	if printstats {
		fmt.Printf("S_%d:", n-1)
		for k := 0; k < n; k++ {
			fmt.Printf(" %d", this.S[k])
		}
		fmt.Printf("\n\tj_%d=%d\n", n-1, this.the_j)
		fmt.Printf("\tj_%d=%d\n", n, this.next_j)
		fmt.Printf("\tS_%d[j_%d]=%d\n", n-1, n-1, this.posify(this.S[this.the_j]))
		if this.S[1] != 0 {
			fmt.Print("\tS[1]!=0")
		}
		fmt.Println()
	}
}

func (this *RC4Cryption) nextVal() int8 {
	this.the_i = (this.the_i + 1) % 256
	this.the_j = (this.the_j + this.posify(this.S[this.the_i])) % 256
	this.sswap(&(this.S), this.the_i, this.the_j)
	pos := (this.posify(this.S[this.the_i]) + this.posify(this.S[this.the_j])) % 256
	value := this.S[pos]
	return value
}
func (this RC4Cryption) posify(byt int8) int {
	if byt >= 0 {
		return int(byt)
	}
	return 256 + int(byt)
}
func (this RC4Cryption) sswap(S *[]int8, i, j int) {
	tmp := (*S)[i]
	(*S)[i] = (*S)[j]
	(*S)[j] = tmp
}
func (this RC4Cryption) ksa(key []byte) {
	this.ksa1(256, key, false)
}
func Encrypt(key []byte, content []byte) []byte {
	enData := make([]byte, len(content))
	rc4 := create()
	rc4.ksa(key)
	rc4.init()
	size := len(content)
	for i := 0; i < size; i++ {
		val := rc4.nextVal()
		enData[i] = byte(int8(content[i]) ^ val)
	}
	return enData
}

func GenerateKeyForRC4(key, id string) []byte {
	value := []byte(id)
	keyBytes, _ := base64.StdEncoding.DecodeString(key)
	buffer := new(bytes.Buffer)
	buffer.Write(keyBytes)
	buffer.WriteByte('_')
	buffer.Write(value)
	return buffer.Bytes()
}
