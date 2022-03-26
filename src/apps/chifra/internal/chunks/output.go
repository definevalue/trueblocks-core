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
		return opts.HandleChunksExtractBlooms()
	} else if opts.Extract == "pins" {
		return opts.HandleChunksExtractPins()
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
		err = opts.HandleChunksExtractBlooms()
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
	}
	// opts.Globals.PassItOn("chunkMan", opts.ToCmdLine())
	return false
	// EXISTING_CODE
}

// EXISTING_CODE
/*
bool COptions::handle_list(void) {
    ASSERT(pins.size());  // local pins have already been loaded
    for (auto pin : pins) {
        if (!isJson()) {
            cout << trim(pin.Format(expContext().fmtMap["format"]), '\t') << endl;
        } else {
            cout << ((isJson() && !firstOut) ? ", " : "");
            indent();
            pin.toJson(cout);
            unindent();
        }
        firstOut = false;
    }
    return false;
}

//----------------------------------------------------------------
bool COptions::hand le_check() {
    bool enabled = getGlobalConfig("chunkMan")->getConfigBool("enabled", "download_manifest", true);
    if (!enabled) {
        LOG_INFO("Manifest not downloaded. Not initializing.");
        return true;
    }

    establishIndexFolders();

    // If the user is calling here, she wants a fresh read even if we've not just freshened.
    pins12.clear();
    pinlib_readManifest(pins12);
    for (auto pin : pins12) {
        string_q source = indexFolder_blooms + pin.fileName + ".bloom";
        copyFile(source, "./thisFile");
        source = "./thisFile";
        string_q cmd1 = "rm -f ./thisFile.gz";  // + " 2>/dev/null";
        if (system(cmd1.c_str())) {
        }                                      // Don't remove cruft. Silences compiler warnings
        cmd1 = "yes | gzip -n -k ./thisFile";  // + " 2>/dev/null";
        if (system(cmd1.c_str())) {
        }                              // Don't remove cruft. Silences compiler warnings
        cmd1 = "ls -l ./thisFile.gz";  // + " 2>/dev/null";
        if (system(cmd1.c_str())) {
        }                                                         // Don't remove cruft. Silences compiler warnings
        cmd1 = "/usr/local/bin/ipfs add thisFile.gz >/tmp/file";  // + " 2>/dev/null";
        if (system(cmd1.c_str())) {
        }  // Don't remove cruft. Silences compiler warnings
        // clang-format on
        LOG_INFO("zip: ", source + ".gz", " ", fileExists(source + ".gz"));
        string_q ret = asciiFileToString("/tmp/file");
        if (ret != pin.bloomHash) {
            cerr << endl;
            cerr << "bloom hashes mismatch for file " << pin.fileName << endl;
            cerr << "\tret: " << ret << endl;
            cerr << "\tpin: " << pin.bloomHash << endl;
            cerr << endl;
        }
    }

    LOG_INFO(bBlue, "Pins were checked.                                           ", cOff);
    return true;  // do not continue
}

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
