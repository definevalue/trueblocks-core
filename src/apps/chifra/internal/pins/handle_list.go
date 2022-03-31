// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package pinsPkg

import (
	"net/http"
	"os"
	"sort"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/pinlib"
)

func (opts *PinsOptions) ListInternal() error {
	manifestData, err := pinlib.ReadLocalManifest()
	if err != nil {
		return err
	}

	sort.Slice(manifestData.Pins, func(i, j int) bool {
		iPin := manifestData.Pins[i]
		jPin := manifestData.Pins[j]
		return iPin.FileName < jPin.FileName
	})

	if opts.Globals.TestMode {
		// Shorten the array for testing
		manifestData.Pins = manifestData.Pins[:100]
	}

	opts.PrintManifestHeader()
	if opts.Globals.ApiMode {
		opts.Globals.Respond(opts.Globals.Writer, http.StatusOK, manifestData.Pins)

	} else {
		err = opts.Globals.Output(os.Stdout, opts.Globals.Format, manifestData.Pins)
		if err != nil {
			logger.Log(logger.Error, err)
		}
	}

	return nil
}
