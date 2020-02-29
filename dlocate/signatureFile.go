package main

import (
	"math/big"
)

const signatureFileCap = 2000
const hierarchicalSignatureFileCap = 100

// SignatureFile contains signature bits that represent the files inside a partition
type SignatureFile struct {
	len  int
	data big.Int
}

func newSignatureFile() SignatureFile {
	return SignatureFile{len: signatureFileCap, data: *big.NewInt(0)}
}

func newSignatureFileH() SignatureFile {
	return SignatureFile{len: hierarchicalSignatureFileCap, data: *big.NewInt(0)}
}

func (sf *SignatureFile) or(sf2 SignatureFile) {
	for i := 0; i < sf2.len; i++ {
		if sf2.getBit(i) {
			sf.setBit(i)
		}
	}
}

func (sf *SignatureFile) setBit(index int) {
	index = index % sf.len
	sf.data.SetBit(&sf.data, index, 1)
}

func (sf *SignatureFile) getBit(index int) bool {
	index = index % sf.len
	return sf.data.Bit(index) == 1
}
