// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package bloomPkg

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
)

func Test_Bloom(t *testing.T) {
	var bloomPath cache.Path
	bloomPath.New("mainnet", cache.BloomChunk)
	path := bloomPath.GetFullPath("000000000-000000000")

	var blooms BloomFilter
	blooms.ReadBloomArray(path)
}
