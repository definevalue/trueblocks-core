// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * This file was auto generated with makeClass --gocmds. DO NOT EDIT.
 */

package cmd

// EXISTING_CODE
import (
	"os"

	chunksPkg "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/internal/chunks"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/internal/globals"
	"github.com/spf13/cobra"
)

// EXISTING_CODE

// chunksCmd represents the chunks command
var chunksCmd = &cobra.Command{
	Use:     usageChunks,
	Short:   shortChunks,
	Long:    longChunks,
	Version: versionText,
	RunE:    chunksPkg.RunChunks,
}

var usageChunks = `chunks [flags] <block> [block...]

Arguments:
  blocks - optional list of blocks to intersect with chunk ranges`

var shortChunks = "manage and investigate chunks and bloom filters"

var longChunks = `Purpose:
  Manage and investigate chunks and bloom filters.`

var notesChunks = `
Notes:
  - Only a single block in a given chunk needs to be supplied.`

func init() {
	chunksCmd.Flags().SortFlags = false

	chunksCmd.Flags().StringVarP(&chunksPkg.GetOptions().Extract, "extract", "e", "", `show some or all of the contents of the chunk or bloom filters
One of [ stats | pins | blooms | index | header | addresses | appearances ]`)
	globals.InitGlobals(chunksCmd, &chunksPkg.GetOptions().Globals)

	chunksCmd.SetUsageTemplate(UsageWithNotes(notesChunks))
	chunksCmd.SetOut(os.Stderr)

	// EXISTING_CODE
	// EXISTING_CODE

	chifraCmd.AddCommand(chunksCmd)
}
