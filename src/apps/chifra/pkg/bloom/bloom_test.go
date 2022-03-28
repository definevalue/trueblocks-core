// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package bloomPkg

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
)

func Test_Bloom(t *testing.T) {
	bloomPath := cache.NewCachePath("mainnet", cache.Index_Bloom)

	var blooms BloomFilter
	blooms.ReadBloomFilter(bloomPath.GetFullPath("000000000-000000000"))
}
