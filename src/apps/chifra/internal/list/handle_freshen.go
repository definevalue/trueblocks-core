package listPkg

// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

// TODO: BOGUS -- USED TO BE ACCTSCRAPE2
import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/index"
)

func (optsEx *MonitorUpdate) HandleFreshenMonitors() error {
	chain := optsEx.Globals.Chain
	indexPath := config.GetPathToIndex(chain) + "finalized/"
	files, err := ioutil.ReadDir(indexPath)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	resultChannel := make(chan []index.ResultRecord, len(files))

	taskCount := 0
	for _, info := range files {
		if !info.IsDir() {
			if taskCount >= optsEx.maxTasks {
				resArray := <-resultChannel
				for _, r := range resArray {
					optsEx.UpdateMonitors(&r)
				}
				taskCount--
			}
			taskCount++
			indexFileName := indexPath + "/" + info.Name()
			wg.Add(1)
			go optsEx.visitChunkToFreshenFinal(indexFileName, resultChannel, &wg)
		}
	}

	wg.Wait()
	close(resultChannel)

	for resArray := range resultChannel {
		for _, r := range resArray {
			optsEx.UpdateMonitors(&r)
		}
	}

	return nil
}

// visitChunkToFreshenFinal opens one index file, searches for all the address(es) we're looking for and pushes the resultRecords
// (even if empty) down the resultsChannel.
func (optsEx *MonitorUpdate) visitChunkToFreshenFinal(fileName string, resultChannel chan<- []index.ResultRecord, wg *sync.WaitGroup) {
	var results []index.ResultRecord
	defer func() {
		wg.Done()
		resultChannel <- results
	}()

	indexChunk, err := index.LoadIndexHeader(fileName)
	if err != nil {
		log.Println(err)
		return
	}
	defer indexChunk.Close()

	if optsEx.Globals.TestMode && indexChunk.Range.Last > 5000000 {
		return
	}

	if !optsEx.RangesIntersect(indexChunk.Range) {
		return
	}

	for _, mon := range optsEx.monMap {
		rec := indexChunk.GetAppearanceRecords(mon.Address)
		if rec != nil {
			results = append(results, *rec)
		} else {
			err := mon.WriteHeader(mon.Deleted, uint32(indexChunk.Range.Last))
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// UpdateMonitors writes an array of appearances to the Monitor file updating the header for lastScanned. It
// is called by 'chifra list' and 'chifra export' prior to reporting results
func (optsEx *MonitorUpdate) UpdateMonitors(result *index.ResultRecord) {
	mon := optsEx.monMap[result.Address]
	if result == nil {
		fmt.Println("Should not happen -- null result")
		return
	}

	if result.AppRecords == nil || len(*result.AppRecords) == 0 {
		fmt.Println("Should not happen -- empty result record:", result.Address)
		return
	}

	_, err := mon.WriteApps(*result.AppRecords, uint32(result.Range.Last))
	if err != nil {
		log.Println(err)
	} else {
		bBlue := (colors.Bright + colors.Blue)
		log.Printf("Found %s%s%s adding appearances (count: %d)\n", bBlue, mon.GetAddrStr(), colors.Off, len(*result.AppRecords))
	}
	// theWriter := csv.NewWriter(opts Ex.writer)
	// theWriter.Comma = 0x9
	// var out [][]string
	// for _, app := range *result.AppRecords {
	// 	out = append(out, []string{strings.ToLower(result.Address.Hex()), fmt.Sprintf("%d", app.BlockNumber), fmt.Sprintf("%d", app.TransactionId)})
	// }
	// theWriter.WriteAll(out)
	// if err := theWriter.Error(); err != nil {
	// 	// TODO: BOGUS
	// 	log.Fatal("F", err)
	// }
}
