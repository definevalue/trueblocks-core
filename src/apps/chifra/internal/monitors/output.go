// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were generated with makeClass --run. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package monitorsPkg

// EXISTING_CODE
import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/monitor"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
	"github.com/spf13/cobra"
)

// EXISTING_CODE

func RunMonitors(cmd *cobra.Command, args []string) error {
	opts := MonitorsFinishParse(args)

	err := opts.ValidateMonitors()
	if err != nil {
		return err
	}

	// EXISTING_CODE
	if opts.HandleCrudCommands() {
		return nil
	}
	return opts.Globals.PassItOn("acctExport", opts.ToCmdLine())
	// EXISTING_CODE
}

func ServeMonitors(w http.ResponseWriter, r *http.Request) bool {
	opts := FromRequest(w, r)

	err := opts.ValidateMonitors()
	if err != nil {
		opts.Globals.RespondWithError(w, http.StatusInternalServerError, err)
		return true
	}

	// EXISTING_CODE
	if !opts.Globals.TestMode { // our test harness does not use DELETE
		delOptions := "--delete, --undelete, or --remove"
		if r.Method == "DELETE" {
			if !opts.Delete && !opts.Undelete && !opts.Remove {
				err = validate.Usage("Specify one of {0} when using the DELETE route.", delOptions)
			}
		} else {
			if opts.Delete || opts.Undelete || opts.Remove {
				delOptions = strings.Replace(delOptions, " or ", " and ", -1)
				err = validate.Usage("The {0} options are not valid when using the GET route.", delOptions)
			}
		}
		if err != nil {
			opts.Globals.RespondWithError(w, http.StatusInternalServerError, err)
			return true
		}
	}
	return opts.HandleCrudCommands()
	// EXISTING_CODE
}

// EXISTING_CODE
func (opts *MonitorsOptions) HandleCrudCommands() bool {
	if !(opts.Delete || opts.Undelete || opts.Remove) {
		return false
	}

	for _, addr := range opts.Addrs {
		m := monitor.NewMonitor(opts.Globals.Chain, addr, false)
		if !file.FileExists(m.Path()) {
			fmt.Println("Monitor not found for address", m.GetAddrStr())
			return true
		} else {
			if opts.Delete {
				m.Delete()
				fmt.Println("Monitor", m.GetAddrStr(), "was deleted but not removed.")
			} else if opts.Undelete {
				m.UnDelete()
				fmt.Println("Monitor", m.GetAddrStr(), "was undeleted.")
			}

			if opts.Remove {
				wasRemoved, err := m.Remove()
				if !wasRemoved || err != nil {
					log.Println("Monitor for ", addr, "was not removed:", err)
					return true
				} else {
					fmt.Println("Monitor for ", addr, "was permanently removed.")
				}
			}
		}
	}
	return true
}

// EXISTING_CODE
