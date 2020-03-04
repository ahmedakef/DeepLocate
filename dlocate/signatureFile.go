package main

import (
	"math/big"
)

const signatureFileCap = 2000
const hierarchicalSignatureFileCap = 100

// SignatureFile contains signature bits that represent the files inside a partition
type SignatureFile struct {
	Len  int
	Data big.Int
}

func newSignatureFile() SignatureFile {
	return SignatureFile{Len: signatureFileCap, Data: *big.NewInt(0)}
}

func newSignatureFileH() SignatureFile {
	return SignatureFile{Len: hierarchicalSignatureFileCap, Data: *big.NewInt(0)}
}

func (sf *SignatureFile) or(sf2 SignatureFile) {
	for i := 0; i < sf2.Len; i++ {
		if sf2.getBit(i) {
			sf.setBit(i)
		}
	}
}

func (sf *SignatureFile) setBit(index int) {
	index = index % sf.Len
	sf.Data.SetBit(&sf.Data, index, 1)
}

func (sf *SignatureFile) getBit(index int) bool {
	index = index % sf.Len
	return sf.Data.Bit(index) == 1
}
