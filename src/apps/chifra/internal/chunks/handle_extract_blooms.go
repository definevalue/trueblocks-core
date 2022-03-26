// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package chunksPkg

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

/*
uint32_t bytesPerLine = (2048 / 64);
//----------------------------------------------------------------
static bool bloomVisitFunc(const string_q& path, void* data) {
    if (endsWith(path, "/")) {
        return forEveryFileInFolder(path + "*", bloomVisitFunc, data);

    } else {
        blknum_t endBlock = NOPOS;
        blknum_t startBlock = path_2_Bn(path, endBlock);
        blknum_t last = *(blknum_t*)data;
        if (last > startBlock)
            return true;

        CBloomArray blooms;
        readBloomFromBinary(path, blooms);

        ostringstream os;
        cout << "range: {" << startBlock << " " << endBlock << "}" << endl;
        cout << "nBlooms: " << blooms.size() << endl;
        cout << "byteWidth: " << getBloomWidthInBytes() << endl;
        for (auto bloom : blooms) {
            cout << "nInserted: " << bloom.nInserted << endl;
            for (size_t i = 0; i < getBloomWidthInBytes(); i++) {
                if (!(i % bytesPerLine)) {
                    if (i != 0)
                        cout << endl;
                    cout << padNum7T(uint64_t(i)) << ": ";
                }
                uint8_t ch = bloom.bits[i];
                cout << bitset<8>(ch) << ' ';
            }
            cout << endl;
            if (isTestMode()) {
                return false;
            }
        }
    }

    return !shouldQuit();
}
*/
