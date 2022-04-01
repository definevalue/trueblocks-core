// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package bloomPkg

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func Test_Bloom(t *testing.T) {
	// bloomPath := cache.NewCachePath("mainnet", cache.Index_Bloom)

	// var blooms BloomFilter
	// blooms.ReadBloomFilter(bloomPath.GetFullPath("000000000-000000000"))

	tests := []struct {
		Addr   common.Address
		Parts  [5]string
		Values [5]uint32
		Bits   [5]uint32
		Member bool
	}{
		{
			Addr:   common.HexToAddress("0x0371a82e4a9d0a4312f3ee2ac9c6958512891372"),
			Parts:  [5]string{"0371a82e", "4a9d0a43", "12f3ee2a", "c9c69585", "12891372"},
			Values: [5]uint32{57780270, 1251805763, 317976106, 3385234821, 310973298},
			Bits:   [5]uint32{108590, 854595, 257578, 431493, 594802},
			Member: true,
		},
		{
			Addr:   common.HexToAddress("0x3d493c51a916f86d6d1c04824b3a7431e61a3ca3"),
			Parts:  [5]string{"3d493c51", "a916f86d", "6d1c0482", "4b3a7431", "e61a3ca3"},
			Values: [5]uint32{1028209745, 2836854893, 1830552706, 1262122033, 3860479139},
			Bits:   [5]uint32{605265, 456813, 787586, 685105, 670883},
			Member: true,
		},
		{
			Addr:   common.HexToAddress("0xe1c15164dcfe79431f8421b5a311a829cf0907f3"),
			Parts:  [5]string{"e1c15164", "dcfe7943", "1f8421b5", "a311a829", "cf0907f3"},
			Values: [5]uint32{3787542884, 3707664707, 528753077, 2735843369, 3473475571},
			Bits:   [5]uint32{86372, 948547, 270773, 108585, 591859},
			Member: true,
		},
		{
			Addr:   common.HexToAddress("0x1296d3bb6dc0efbae431a12939fc15b2823db49b"),
			Parts:  [5]string{"1296d3bb", "6dc0efba", "e431a129", "39fc15b2", "823db49b"},
			Values: [5]uint32{311874491, 1841360826, 3828457769, 972821938, 2185082011},
			Bits:   [5]uint32{447419, 61370, 106793, 791986, 898203},
			Member: true,
		},
		{
			Addr:   common.HexToAddress("0xd09022c48298f268c2c431dadb9ca4c2534d9c1c"),
			Parts:  [5]string{"d09022c4", "8298f268", "c2c431da", "db9ca4c2", "534d9c1c"},
			Values: [5]uint32{3499107012, 2191061608, 3267637722, 3684476098, 1397595164},
			Bits:   [5]uint32{8900, 586344, 274906, 828610, 891932},
			Member: true,
		},
		{
			Addr:   common.HexToAddress("0x1296d3bb6dc0efbae431a12939fc15b2823db79b"),
			Parts:  [5]string{"1296d3bb", "6dc0efba", "e431a129", "39fc15b2", "823db79b"},
			Values: [5]uint32{311874491, 1841360826, 3828457769, 972821938, 2185082779},
			Bits:   [5]uint32{447419, 61370, 106793, 791986, 898971},
			Member: false,
		},
		{
			Addr:   common.HexToAddress("0xd09022c48298f268c2c431dadb9ca4c2534d9c1e"),
			Parts:  [5]string{"d09022c4", "8298f268", "c2c431da", "db9ca4c2", "534d9c1e"},
			Values: [5]uint32{3499107012, 2191061608, 3267637722, 3684476098, 1397595166},
			Bits:   [5]uint32{8900, 586344, 274906, 828610, 891934},
			Member: false,
		},
		{
			Addr:   common.HexToAddress("0x0341a82e4a9d0a4312f3ee2ac9c6958512891342"),
			Parts:  [5]string{"0341a82e", "4a9d0a43", "12f3ee2a", "c9c69585", "12891342"},
			Values: [5]uint32{54634542, 1251805763, 317976106, 3385234821, 310973250},
			Bits:   [5]uint32{108590, 854595, 257578, 431493, 594754},
			Member: false,
		},
		{
			Addr:   common.HexToAddress("0x3d493c51a916f86d6d1c04824b3a4431e61a3ca3"),
			Parts:  [5]string{"3d493c51", "a916f86d", "6d1c0482", "4b3a4431", "e61a3ca3"},
			Values: [5]uint32{1028209745, 2836854893, 1830552706, 1262109745, 3860479139},
			Bits:   [5]uint32{605265, 456813, 787586, 672817, 670883},
			Member: false,
		},
		{
			Addr:   common.HexToAddress("0xe1c15164dcfe49431f8421b5a311a829cf0904f3"),
			Parts:  [5]string{"e1c15164", "dcfe4943", "1f8421b5", "a311a829", "cf0904f3"},
			Values: [5]uint32{3787542884, 3707652419, 528753077, 2735843369, 3473474803},
			Bits:   [5]uint32{86372, 936259, 270773, 108585, 591091},
			Member: false,
		},
	}

	bloom := NewBloomFilter()
	for _, tt := range tests {
		_, _, bits := bitsToLight(tt.Addr)
		// for i := 0; i < 5; i++ {
		// 	// fmt.Println(tt.Addr, ":")
		// 	// b := []byte(tt.Parts[i])
		// 	// fmt.Printf("\tParts: %v %v\n", b, parts[i])
		// 	// if tt.Values[i] == values[i] {
		// 	// 	fmt.Println("\tValues:", tt.Values[i], values[i], (tt.Values[i] == values[i]))
		// 	// }
		// 	// if tt.Bits[i] == bits[i] {
		// 	// 	fmt.Println("\tBits:", tt.Bits[i], bits[i], (tt.Bits[i] == bits[i]))
		// 	// }
		// }
		if tt.Member {
			bloom.LightBits(bits)
		}
	}

	expectedLit := []uint64{
		8900, 61370, 86372, 106793, 108585, 108590, 257578,
		270773, 274906, 431493, 447419, 456813, 586344, 591859,
		594802, 605265, 670883, 685105, 787586, 791986, 828610,
		854595, 891932, 898203, 948547,
	}

	nBlooms, nInserted, nBitsLit, nBitsNotLit, sz, bitsLit := bloom.getStats()
	fmt.Println(nBlooms, nInserted, nBitsLit, nBitsNotLit, sz, bitsLit)

	if len(bitsLit) != len(expectedLit) {
		t.Error("mismatched lengths -- expected:", len(expectedLit), "got:", len(bitsLit))
	} else {
		for i, bit := range bitsLit {
			if bit != expectedLit[i] {
				t.Error("mismatched bit lit -- expected:", expectedLit[i], "got:", bitsLit[i])
			}
		}
	}
	for _, tt := range tests {
		if tt.Member && !bloom.isInBloomFilter(tt.Addr) {
			t.Error("Address should be member, but isn't", tt.Addr.Hex())
		} else if !tt.Member && bloom.isInBloomFilter(tt.Addr) { // && !tt.FalsePositive {
			t.Error("Address should not be member, but is (ignores false positives)", tt.Addr.Hex())
		}
		fmt.Println(strings.ToLower(tt.Addr.Hex()), bloom.isInBloomFilter(tt.Addr))
	}

	// t.Error("what")
}

// writeBloom
//     lockSection();
//     CArchive output(WRITING_ARCHIVE);
//     if (!output.Lock(fileName, modeWriteCreate, LOCK_NOWAIT)) {
//         unlockSection();
//         return false;
//     }
//     output.Write((uint32_t)blooms.size());
//     for (auto bloom : blooms) {
//         output.Write(bloom.nInserted);
//         output.Write(bloom.bits, sizeof(uint8_t), BLOOM_WIDTH_IN_BYTES);
//     }
//     output.Release();
//     unlockSection();
//     return true;
// }
