package bloomPkg

// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
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
	BLOOM_WIDTH_IN_BYTES = (1048576 / 8)
)

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

const bytesPerLine = (2048 / 64)

func (bloom *BloomFilter) DebugPrint() {
	fmt.Println("range:", bloom.Range)
	fmt.Println("nBlooms:", bloom.Count)
	fmt.Println("byteWidth:", BLOOM_WIDTH_IN_BYTES)
	for i := uint32(0); i < bloom.Count; i++ {
		fmt.Println("nInserted:", bloom.Blooms[i].NInserted)
		for j := 0; j < len(bloom.Blooms[i].Bytes); j++ {
			if (j % bytesPerLine) == 0 {
				if j != 0 {
					fmt.Println()
				}
				fmt.Printf("% 7d: ", j)
			}
			ch := bloom.Blooms[i].Bytes[j]
			fmt.Printf("%08b ", ch)
		}
		fmt.Println()
	}
}
