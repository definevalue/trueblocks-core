// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were generated with makeClass --run. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package chunksPkg

// EXISTING_CODE
import (
	"net/http"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/spf13/cobra"
)

// EXISTING_CODE

func RunChunks(cmd *cobra.Command, args []string) error {
	opts := ChunksFinishParse(args)

	err := opts.ValidateChunks()
	if err != nil {
		return err
	}

	// EXISTING_CODE
	if opts.Extract == "blooms" {
		return opts.HandleChunksExtract(opts.showBloom)
	} else if opts.Extract == "pins" {
		return opts.HandleChunksExtractPins()
	} else if opts.Extract == "stats" {
		return opts.HandleChunksExtract(opts.showStats)
	}

	return opts.Globals.PassItOn("chunkMan", opts.ToCmdLine())
	// EXISTING_CODE
}

func ServeChunks(w http.ResponseWriter, r *http.Request) bool {
	opts := FromRequest(w, r)

	err := opts.ValidateChunks()
	if err != nil {
		opts.Globals.RespondWithError(w, http.StatusInternalServerError, err)
		return true
	}

	// EXISTING_CODE
	if opts.Extract == "blooms" {
		err = opts.HandleChunksExtract(opts.showBloom)
		if err != nil {
			logger.Log(logger.Warning, "Could not extract blooms", err)
		}
		return true
	} else if opts.Extract == "pins" {
		err = opts.HandleChunksExtractPins()
		if err != nil {
			logger.Log(logger.Warning, "Could not extract pin list", err)
		}
		return true
	} else if opts.Extract == "stats" {
		err = opts.HandleChunksExtract(opts.showStats)
		if err != nil {
			logger.Log(logger.Warning, "Could not extract stats", err)
		}
		return true
	}
	// opts.Globals.PassItOn("chunkMan", opts.ToCmdLine())
	return false
	// EXISTING_CODE
}

// EXISTING_CODE
/*
if (share) {
	    string_q res := doCommand("which ipfs");
	    if (res.empty()) {
	        return usa ge("Could not find ipfs in your $PATH. You must install ipfs for the --share command to work.");
		}
	}
	if (share) {
	    ostringstream os;
	    os << "ipfs add -Q --pin \"" << bloomFn + "\"";
	    string_q newHash = doCommand(os.str());
	    LOG_INFO(cGreen, "Re-pinning ", pin.fileName, cOff, " ==> ", newHash, " ",
	         (pin.bloomHash == newHash ? greenCheck : redX));
	}

*/
// EXISTING_CODE
