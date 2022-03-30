package bloomPkg

// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/ethereum/go-ethereum/common"
)

type bloomBytes struct {
	NInserted uint32
	Bytes     []byte
}

type BloomFilter struct {
	Range  cache.FileRange
	Count  uint32
	Blooms []bloomBytes
}

const (
	BLOOM_WIDTH_IN_BITS  = (1048576)
	BLOOM_WIDTH_IN_BYTES = (BLOOM_WIDTH_IN_BITS / 8)
)

func NewBloomFilter() BloomFilter {
	var ret BloomFilter
	ret.Blooms = make([]bloomBytes, 1)
	ret.Count = 1
	return ret
}

func NewBloomFilterFromAddress(addr common.Address) BloomFilter {
	ret := NewBloomFilter()
	ret.addAddrToBloom(addr)
	return ret
}

func chunkBytes(addr common.Address, chunkSize int) [][]byte {
	var chunks [][]byte
	slice := addr.Bytes()
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// bloom_t addr_2_Bloom2(const address_t& addr, CUintArray& litBits) {
//     bloom_t ret;
//     cout << addr << " ";
//     for (size_t k = 0; k < K; k++) {
//         string_q four_byte = extract(addr, 2 + (k * NIBBLE_WID), NIBBLE_WID);
//         uint64_t bit64 = str_2_Uint("0x" + four_byte);
//         uint64_t bit = (bit64 % BLOOM_WIDTH_IN_BITS);
//         ret.lightBit(bit);
//         litBits.push_back(bit);
//         cout << four_byte << "-" << uint_2_Str(bit64) << "-" << uint_2_Str(bit) << " ";
//     }
//     cout << endl;
//     return ret;
// }
func (bloom *BloomFilter) addAddrToBloom(addr common.Address) error {
	return nil
}

func (bloom *BloomFilter) ReadBloomFilter(fileName string) (err error) {
	bloom.Range, err = cache.RangeFromFilename(fileName)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	err = binary.Read(file, binary.LittleEndian, &bloom.Count)
	if err != nil {
		return err
	}

	bloom.Blooms = make([]bloomBytes, bloom.Count)
	for i := uint32(0); i < bloom.Count; i++ {
		// fmt.Println("nBlooms:", bloom.Count)
		err = binary.Read(file, binary.LittleEndian, &bloom.Blooms[i].NInserted)
		if err != nil {
			return err
		}
		// fmt.Println("nInserted:", bloom.Blooms[i].NInserted)
		bloom.Blooms[i].Bytes = make([]byte, BLOOM_WIDTH_IN_BYTES)
		err = binary.Read(file, binary.LittleEndian, &bloom.Blooms[i].Bytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func (bloom *BloomFilter) DisplayBloom(verbose int) {
	var bytesPerLine = (2048 / 16)
	if verbose > 0 {
		if verbose > 4 {
			bytesPerLine = 128
		} else {
			bytesPerLine = 32
		}
	}

	nInserted := uint32(0)
	for i := uint32(0); i < bloom.Count; i++ {
		nInserted += bloom.Blooms[i].NInserted
	}
	fmt.Println("range:", bloom.Range)
	fmt.Println("nBlooms:", bloom.Count)
	fmt.Println("byteWidth:", BLOOM_WIDTH_IN_BYTES)
	fmt.Println("nInserted:", nInserted)
	if verbose > 0 {
		for i := uint32(0); i < bloom.Count; i++ {
			for j := 0; j < len(bloom.Blooms[i].Bytes); j++ {
				if (j % bytesPerLine) == 0 {
					if j != 0 {
						fmt.Println()
					}
				}
				ch := bloom.Blooms[i].Bytes[j]
				str := fmt.Sprintf("%08b", ch)
				fmt.Printf("%s ", strings.Replace(str, "0", ".", -1))
			}
		}
		fmt.Println()
	}
}
