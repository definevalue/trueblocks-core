// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package chunksPkg

import (
	bloomPkg "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/bloom"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
)

func (opts *ChunksOptions) HandleChunksExtractBlooms() error {
	bloomPath := cache.NewCachePath(opts.Globals.Chain, cache.Index_Bloom)
	path := bloomPath.GetFullPath("000000000-000000000")

	var bloom bloomPkg.BloomFilter
	bloom.ReadBloomFilter(path)
	bloom.DisplayBloom(int(opts.Globals.LogLevel))

	bloomPath = cache.NewCachePath(opts.Globals.Chain, cache.Index_Bloom)
	path = bloomPath.GetFullPath("000000001-000590501")
	bloomArray.ReadBloomFilter(path)
	bloomArray.DisplayBloom(1)

	return nil
}
