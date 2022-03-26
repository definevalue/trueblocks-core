// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package chunksPkg

import (
	"net/http"
	"os"
	"sort"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/pinlib/manifest"
)

func (opts *ChunksOptions) HandleChunksExtractPins() error {
	manifestData, err := manifest.FromLocalFile(opts.Globals.Chain)
	if err != nil {
		return err
	}

	sort.Slice(manifestData.NewPins, func(i, j int) bool {
		iPin := manifestData.NewPins[i]
		jPin := manifestData.NewPins[j]
		return iPin.FileName < jPin.FileName
	})

	if opts.Globals.TestMode {
		// Shorten the array for testing
		manifestData.NewPins = manifestData.NewPins[:100]
	}

	// TODO: Fix export without arrays
	opts.PrintManifestHeader()
	if opts.Globals.ApiMode {
		opts.Globals.Respond(opts.Globals.Writer, http.StatusOK, manifestData.NewPins)

	} else {
		err = opts.Globals.Output(os.Stdout, opts.Globals.Format, manifestData.NewPins)
		if err != nil {
			logger.Log(logger.Error, err)
		}
	}

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

func (opts *ChunksOptions) PrintManifestHeader() {
	// The following two values should be read the manifest, however right now only
	// TSV format is available for download and it lacks this information
	// TODO: These values should be in a config file
	// TODO: We can add the "loaded" configuration file to Options
	// TODO: This needs to be per chain data
	chain := opts.Globals.Chain
	logger.Log(logger.Info, "hashToIndexFormatFile:", "Qmart6XP9XjL43p72PGR93QKytbK8jWWcMguhFgxATTya2")
	logger.Log(logger.Info, "hashToBloomFormatFile:", "QmNhPk39DUFoEdhUmtGARqiFECUHeghyeryxZM9kyRxzHD")
	logger.Log(logger.Info, "manifestHashEncoding:", config.ReadBlockScrape(chain).UnchainedIndex.ManifestHashEncoding)
	logger.Log(logger.Info, "unchainedIndexAddr:", config.ReadBlockScrape(chain).UnchainedIndex.Address)
	if !opts.Globals.TestMode {
		logger.Log(logger.Info, "manifestLocation:", config.GetPathToChainConfig(chain)) // order matters
		logger.Log(logger.Info, "unchainedIndexFolder:", config.GetPathToIndex(chain))   // order matters
	}
}
