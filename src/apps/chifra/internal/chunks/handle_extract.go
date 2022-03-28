// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package chunksPkg

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
)

func (opts *ChunksOptions) HandleChunksExtract(displayFunc func(path string)) error {
	blocks := cache.Convert(opts.Blocks)
	filenameChan := make(chan cache.IndexFileInfo)

	go cache.WalkCacheFolder(opts.Globals.Chain, cache.Index_Bloom, filenameChan)

	for result := range filenameChan {
		switch result.Type {
		case cache.Index_Bloom:
			hit := false
			for _, block := range blocks {
				h := cache.BlockIntersects(result.Range, block)
				hit = hit || h
			}
			if len(blocks) == 0 || hit {
				displayFunc(result.Path)
			}
		case cache.None:
			close(filenameChan)
		}
	}

	return nil
}
