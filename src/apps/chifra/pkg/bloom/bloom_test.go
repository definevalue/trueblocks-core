// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package bloomPkg

import (
	"encoding/binary"
	"fmt"
	"strings"
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/ethereum/go-ethereum/common"
)

func Test_Bloom(t *testing.T) {
	bloomPath := cache.NewCachePath("mainnet", cache.Index_Bloom)

	var blooms BloomFilter
	blooms.ReadBloomFilter(bloomPath.GetFullPath("000000000-000000000"))

	tests := []struct {
		Addr      common.Address
		Expected1 string
	}{
		{
			Addr:      common.HexToAddress("0x0371a82e4a9d0a4312f3ee2ac9c6958512891372"),
			Expected1: `0371a82e-57780270-108590 4a9d0a43-1251805763-854595 12f3ee2a-317976106-257578 c9c69585-3385234821-431493 12891372-310973298-594802 `,
		},
		{
			Addr:      common.HexToAddress("0x1296d3bb6dc0efbae431a12939fc15b2823db79b"),
			Expected1: `1296d3bb-311874491-447419 6dc0efba-1841360826-61370 e431a129-3828457769-106793 39fc15b2-972821938-791986 823db79b-2185082779-898971 `,
		},
		{
			Addr:      common.HexToAddress("0x3d493c51a916f86d6d1c04824b3a7431e61a3ca3"),
			Expected1: `3d493c51-1028209745-605265 a916f86d-2836854893-456813 6d1c0482-1830552706-787586 4b3a7431-1262122033-685105 e61a3ca3-3860479139-670883 `,
		},
		{
			Addr:      common.HexToAddress("0xd09022c48298f268c2c431dadb9ca4c2534d9c1c"),
			Expected1: `d09022c4-3499107012-8900 8298f268-2191061608-586344 c2c431da-3267637722-274906 db9ca4c2-3684476098-828610 534d9c1c-1397595164-891932 `,
		},
		{
			Addr:      common.HexToAddress("0xe1c15164dcfe79431f8421b5a311a829cf0907f3"),
			Expected1: `e1c15164-3787542884-86372 dcfe7943-3707664707-948547 1f8421b5-528753077-270773 a311a829-2735843369-108585 cf0907f3-3473475571-591859 `,
		},
		{
			Addr:      common.HexToAddress("0x0341a82e4a9d0a4312f3ee2ac9c6958512891342"),
			Expected1: `0341a82e-54634542-108590 4a9d0a43-1251805763-854595 12f3ee2a-317976106-257578 c9c69585-3385234821-431493 12891342-310973250-594754 `,
		},
		{
			Addr:      common.HexToAddress("0x1296d3bb6dc0efbae431a12939fc15b2823db49b"),
			Expected1: `1296d3bb-311874491-447419 6dc0efba-1841360826-61370 e431a129-3828457769-106793 39fc15b2-972821938-791986 823db49b-2185082011-898203 `,
		},
		{
			Addr:      common.HexToAddress("0x3d493c51a916f86d6d1c04824b3a4431e61a3ca3"),
			Expected1: `3d493c51-1028209745-605265 a916f86d-2836854893-456813 6d1c0482-1830552706-787586 4b3a4431-1262109745-672817 e61a3ca3-3860479139-670883 `,
		},
		{
			Addr:      common.HexToAddress("0xd09022c48298f268c2c431dadb9ca4c2534d9c1c"),
			Expected1: `d09022c4-3499107012-8900 8298f268-2191061608-586344 c2c431da-3267637722-274906 db9ca4c2-3684476098-828610 534d9c1c-1397595164-891932 `,
		},
		{
			Addr:      common.HexToAddress("0xe1c15164dcfe49431f8421b5a311a829cf0904f3"),
			Expected1: `e1c15164-3787542884-86372 dcfe4943-3707652419-936259 1f8421b5-528753077-270773 a311a829-2735843369-108585 cf0904f3-3473474803-591091 `,
		},
	}

	for _, tt := range tests {
		got := ""
		fourBytes := chunkBytes(tt.Addr, 4)
		for _, fourByte := range fourBytes {
			fourByteAsUint32 := binary.BigEndian.Uint32(fourByte)
			widthInBits := uint32(BLOOM_WIDTH_IN_BITS)
			bitToLight := (fourByteAsUint32 % widthInBits)
			got += fmt.Sprintf("%x-%d-%d ", fourByte, fourByteAsUint32, bitToLight)
		}
		fmt.Println(strings.ToLower(tt.Addr.Hex()), got)

		if got != tt.Expected1 {
			t.Error("Expected:\n" + tt.Expected1 + "\ngot:\n" + got)
		}
	}

	// t.Error("test")
}

/*

00015430 (  15430)- <INFO>  : 0x0371a82e4a9d0a4312f3ee2ac9c6958512891372-1
00016597 (   1167)- <INFO>  : 0x1296d3bb6dc0efbae431a12939fc15b2823db79b-0
00018020 (   1423)- <INFO>  : 0x3d493c51a916f86d6d1c04824b3a7431e61a3ca3-1
00019543 (   1523)- <INFO>  : 0xd09022c48298f268c2c431dadb9ca4c2534d9c1c-1-fp
00021157 (   1614)- <INFO>  : 0xe1c15164dcfe79431f8421b5a311a829cf0907f3-1
00022042 (    885)- <INFO>  : 0x0341a82e4a9d0a4312f3ee2ac9c6958512891342-0
00023431 (   1389)- <INFO>  : 0x1296d3bb6dc0efbae431a12939fc15b2823db49b-1
00024560 (   1129)- <INFO>  : 0x3d493c51a916f86d6d1c04824b3a4431e61a3ca3-0
00026371 (   1811)- <INFO>  : 0xd09022c48298f268c2c431dadb9ca4c2534d9c1c-1
00027151 (    780)- <INFO>  : 0xe1c15164dcfe49431f8421b5a311a829cf0904f3-0

*/
