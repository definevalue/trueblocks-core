package chunksPkg

// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

import (
	bloomPkg "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/bloom"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
)

func (opts *ChunksOptions) HandleChunksExtractBlooms() error {
	var bloomPath cache.Path
	bloomPath.New(opts.Globals.Chain, cache.BloomChunk)
	path := bloomPath.GetFullPath("000000000-000000000")

	var bloomArray bloomPkg.BloomFilter
	bloomArray.ReadBloomArray(path)
	bloomArray.DebugPrint()

	bloomPath.New(opts.Globals.Chain, cache.BloomChunk)
	path = bloomPath.GetFullPath("000000001-000590501")
	bloomArray.ReadBloomArray(path)
	bloomArray.DebugPrint()

	return nil
}
