// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package listPkg

import (
	"net/http"
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/monitor"
)

func (opts *ListOptions) HandleListCount(results []monitor.Monitor) error {
	// TODO: Fix export without arrays
	if opts.Globals.ApiMode {
		opts.Globals.Respond(opts.Globals.Writer, http.StatusOK, results)

	} else {
		err := opts.Globals.Output(os.Stdout, opts.Globals.Format, results)
		if err != nil {
			logger.Log(logger.Error, err)
		}
	}

	return nil
}
