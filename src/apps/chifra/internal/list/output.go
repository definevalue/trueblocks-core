// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were generated with makeClass --run. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package listPkg

// EXISTING_CODE
import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	exportPkg "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/internal/export"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/internal/globals"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/blockRange"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/monitor"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// AddressMonitorMap carries arrays of appearances that have not yet been written to the monitor file
type AddressMonitorMap map[common.Address]*monitor.Monitor

// MonitorUpdate stores the original 'chifra list' command line options plus
type MonitorUpdate struct {
	writer    io.Writer
	maxTasks  int
	monArray2 []monitor.Monitor
	monMap    AddressMonitorMap
	Globals   globals.GlobalOptions
	Range     blockRange.FileRange
}

// EXISTING_CODE

func RunList(cmd *cobra.Command, args []string) error {
	opts := ListFinishParse(args)

	err := opts.ValidateList()
	if err != nil {
		return err
	}

	// EXISTING_CODE
	if opts.Newone {
		var optsEx MonitorUpdate
		optsEx.writer = os.Stdout
		optsEx.maxTasks = 12
		optsEx.monArray2 = make([]monitor.Monitor, 0, len(opts.Addrs))
		optsEx.monMap = make(AddressMonitorMap, len(opts.Addrs))
		optsEx.Range = blockRange.FileRange{First: 0, Last: utils.NOPOS}
		optsEx.Globals = opts.Globals
		for _, addr := range opts.Addrs {
			if optsEx.monMap[common.HexToAddress(addr)] == nil {
				m := monitor.NewStagedMonitor(optsEx.Globals.Chain, addr)
				optsEx.monArray2 = append(optsEx.monArray2, m)
				optsEx.monMap[m.Address] = &optsEx.monArray2[len(optsEx.monArray2)-1]
			}
			fmt.Println("len:", len(optsEx.monArray2))
			fmt.Println("len:", len(optsEx.monMap))
		}
		err = optsEx.HandleFreshenMonitors()
		if err != nil {
			return err
		}
		optsEx.MoveAllToProduction()
		if opts.Count {
			fmt.Println("len2:", len(optsEx.monArray2))
			fmt.Println("len2:", len(optsEx.monMap))
			return opts.HandleListCount(optsEx.monArray2)
		}
		return nil
	}

	// exportPkg "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/internal/export"
	exportPkg.GetOptions().Appearances = true
	if opts.Count {
		exportPkg.GetOptions().Count = true
	}
	if opts.FirstBlock > 0 {
		exportPkg.GetOptions().FirstBlock = opts.FirstBlock
	}
	if opts.LastBlock > 0 {
		exportPkg.GetOptions().LastBlock = opts.LastBlock
	}
	exportPkg.GetOptions().Globals = opts.Globals
	return exportPkg.RunExport(cmd, args)
	// EXISTING_CODE
}

func ServeList(w http.ResponseWriter, r *http.Request) bool {
	opts := FromRequest(w, r)

	err := opts.ValidateList()
	if err != nil {
		opts.Globals.RespondWithError(w, http.StatusInternalServerError, err)
		return true
	}

	// EXISTING_CODE
	// TODO: BOGUS -- HANDLE THIS IN GOLANG
	return false
	// EXISTING_CODE
}

// EXISTING_CODE
func (optsEx *MonitorUpdate) RangesIntersect(r2 blockRange.FileRange) bool {
	// fmt.Println(r1.First, r1.Last, "-", r2.First, r2.Last, !(r1.Last < r2.First || r1.First > r2.Last))
	return !(optsEx.Range.Last < r2.First || optsEx.Range.First > r2.Last)
}

func (optsEx *MonitorUpdate) MoveAllToProduction() {
	for _, mon := range optsEx.monMap {
		err := mon.MoveToProduction()
		if err != nil {
			log.Println(err)
		}
	}
}

// EXISTING_CODE
