// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package monitor

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/index"
	"github.com/ethereum/go-ethereum/common"
)

func Test_Monitor(t *testing.T) {
	testAddr := "0xF503017d7bAf7fbc0fff7492b751025c6a78179b"

	mon := NewMonitor("mainnet", testAddr)
	path := mon.Path()
	dir, fileName := filepath.Split(path)

	if !strings.HasSuffix(dir, "/cache/mainnet/monitors/") {
		t.Error("Incorrect suffix in 'dir'. Expected: \"/cache/mainnet/monitors/\" Dir:", dir)
	}

	if testAddr+".acct.bin" == fileName {
		t.Error("Filename should be lower case: ", fileName)
	}

	if strings.ToLower(testAddr+".acct.bin") != fileName {
		t.Error("Unexpected filename: ", fileName)
	}
}

func Test_Monitor_Peek(t *testing.T) {
	mon := GetTestMonitor(t)
	defer func() {
		mon.Delete()
		RemoveMonitor(&mon, t)
	}()

	got, err := mon.Peek(0)
	if err == nil {
		t.Error("Should have been 'index out of range in Peek[0]' error")
	}

	expected := index.AppearanceRecord{BlockNumber: 1001001, TransactionId: 1001001}
	got, err = mon.Peek(1)
	if got != expected || err != nil {
		t.Error("Expected:", expected, "Got:", got, err)
	}

	expected = index.AppearanceRecord{BlockNumber: 1001002, TransactionId: 1001002}
	got, err = mon.Peek(2)
	if got != expected || err != nil {
		t.Error("Expected:", expected, "Got:", got, err)
	}

	got, err = mon.Peek(mon.Count)
	if got != expected || err != nil {
		t.Error("Expected:", expected, "Got:", got, err)
	}

	got, err = mon.Peek(3)
	if err == nil {
		t.Error("Should have been 'index out of range in Peek[3]' error")
	}
}

func Test_Monitor_Delete(t *testing.T) {
	mon := GetTestMonitor(t)
	defer func() {
		RemoveMonitor(&mon, t)
	}()

	// The monitor should report that it has two appearances
	got := fmt.Sprintln(mon.ToJSON())
	expected := "{\"address\":\"0xf503017d7baf7fbc0fff7492b751025c6a781791\",\"count\":2,\"fileSize\":16}\n"
	if got != expected {
		t.Error("Expected:", expected, "Got:", got)
	}

	// Try to remove the monitor. It should not be removed because it is not deleted first
	removed, err := mon.Remove()
	if err == nil || removed {
		t.Error("Should not be able to remove monitor without deleting it first")
	} else {
		t.Log("Correctly errors with:", err)
	}

	wasDeleted := mon.ToggleDelete()
	t.Log(mon.ToJSON())
	if wasDeleted || !mon.Deleted {
		t.Error("Should not have been previously deleted, but it should be deleted now")
	}

	wasDeleted = mon.Delete()
	t.Log(mon.ToJSON())
	if !wasDeleted || !mon.Deleted {
		t.Error("Should have been previously deleted, and it should be deleted now")
	}

	wasDeleted = mon.UnDelete()
	t.Log(mon.ToJSON())
	if !wasDeleted || mon.Deleted {
		t.Error("Should have been previously deleted, but should no longer be")
	}

	wasDeleted = mon.Delete()
	t.Log(mon.ToJSON())
	if wasDeleted || !mon.Deleted {
		t.Error("Should not have been previously deleted, but it should be deleted now")
	}
}

func Test_Monitor_Print(t *testing.T) {
	mon := GetTestMonitor(t)
	defer func() {
		mon.Delete()
		RemoveMonitor(&mon, t)
	}()

	// The monitor should report that it has two appearances
	got := fmt.Sprintln(mon.ToJSON())
	expected := "{\"address\":\"0xf503017d7baf7fbc0fff7492b751025c6a781791\",\"count\":2,\"fileSize\":16}\n"
	if got != expected {
		t.Error("Expected:", expected, "Got:", got)
	}
}

func GetTestMonitor(t *testing.T) Monitor {
	// Create a new, empty monitor
	testAddr := "0xF503017d7bAf7fbc0fff7492b751025c6a781791"
	mon := NewMonitor("mainnet", testAddr)

	if mon.Address != common.HexToAddress(testAddr) {
		t.Error("Expected:", common.HexToAddress(testAddr), "Got:", mon.Address)
	}

	if mon.GetAddrStr() != strings.ToLower(testAddr) {
		t.Error("Expected:", strings.ToLower(testAddr), "Got:", mon.GetAddrStr())
	}

	// The file should exist...
	if !file.FileExists(mon.Path()) {
		t.Error("File", mon.Path(), "should exist")
	}

	// and be empty
	if mon.Count != 0 {
		t.Error("New monitor file should be empty")
	}

	apps := []index.AppearanceRecord{
		{BlockNumber: 1001001, TransactionId: 0},
		{BlockNumber: 1001002, TransactionId: 1},
	}
	if len(apps) != 2 {
		t.Error("Incorrect length for test data:", len(apps), "should be 2.")
	}

	// Append two appearances to the monitor
	count, err := mon.Append(apps)
	if err != nil {
		t.Error(err)
	}
	if count != 2 {
		t.Error("Expected count 2 for monitor, got:", count)
	}

	return mon
}

func RemoveMonitor(mon *Monitor, t *testing.T) {
	if !file.FileExists(mon.Path()) {
		t.Error("Monitor file should exist")
	}
	if !mon.Deleted {
		t.Error("Monitor should be deleted")
	}
	if mon.Count != 2 {
		t.Error("Monitor should have two records, has:", mon.Count)
	}
	removed, err := mon.Remove()
	if !removed || err != nil {
		t.Error("Monitor should have been removed", err)
	}
	if file.FileExists(mon.Path()) {
		t.Error("Monitor file should not exist, but it does")
	}
}
