package monitor

// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/index"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// MonitorHeader is the header of the Monitor file
type MonitorHeader struct {
	Deleted     bool   `json:"deleted,omitempty"`
	LastScanned uint32 `json:"lastScanned,omitempty"`
}

// Monitor carries information about a Monitor file and its header
type Monitor struct {
	Address  common.Address `json:"address"`
	Count    uint32         `json:"count"`
	FileSize uint32         `json:"fileSize"`
	Chain    string         `json:"-"`
	ReadFp   *os.File       `json:"-"`
	MonitorHeader
}

// NewMonitor returns a Monitor (but has not yet read in the AppearanceRecords)
func NewMonitor(chain, addr string) Monitor {
	mon := new(Monitor)
	mon.Address = common.HexToAddress(strings.ToLower(addr))
	mon.Chain = chain
	mon.Reload()
	return *mon
}

// String implements the Stringer interface
func (mon Monitor) String() string {
	if mon.Deleted {
		return fmt.Sprintf("%s\t%d\t%d\t%t", hexutil.Encode(mon.Address.Bytes()), mon.Count, mon.FileSize, mon.Deleted)

	}
	return fmt.Sprintf("%s\t%d\t%d", hexutil.Encode(mon.Address.Bytes()), mon.Count, mon.FileSize)
}

func (mon Monitor) ToJSON() string {
	bytes, err := json.Marshal(mon)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func (mon *Monitor) Path() (path string) {
	path = config.GetPathToCache(mon.Chain) + "monitors/" + strings.ToLower(mon.Address.Hex()) + ".acct.bin"
	return
}

func (mon *Monitor) Reload() (uint32, error) {
	path := mon.Path()
	if !file.FileExists(path) {
		// Make sure the file exists since we've been told to monitor it
		file.Touch(path)
	}
	mon.FileSize = uint32(file.FileSize(path))
	mon.Count = uint32(file.FileSize(path) / index.AppRecordWidth)
	return mon.Count, nil
}

func (mon *Monitor) GetAddrStr() string {
	return strings.ToLower(mon.Address.Hex())
}

// ReadApp returns the appearance at the index - 1. For example, ask for idx == 1 to get the
// first record in the file or idx == Count to get the last record in the file.
func (mon *Monitor) ReadApp(idx uint32, app *index.AppearanceRecord) (err error) {
	if idx == 0 || idx > mon.Count {
		// one-based index for ease in caller code
		err = errors.New(fmt.Sprintf("index out of range in ReadApp[%d]", idx))
		return
	}

	if mon.ReadFp == nil {
		path := mon.Path()
		mon.ReadFp, err = os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			return
		}
	}

	// Caller wants record 1, which stands at location 0, etc.
	byteIndex := int64(idx-1) * index.AppRecordWidth
	_, err = mon.ReadFp.Seek(byteIndex, io.SeekStart)
	if err != nil {
		return
	}

	err = binary.Read(mon.ReadFp, binary.LittleEndian, &app.BlockNumber)
	if err != nil {
		return
	}
	err = binary.Read(mon.ReadFp, binary.LittleEndian, &app.TransactionId)
	return
}

// ReadApps returns appearances starting at the first appearance in the file. Use
// make([]index.AppearanceRecord, mon.Count) to create an array big enough
func (mon *Monitor) ReadApps(apps *[]index.AppearanceRecord) (err error) {
	if mon.ReadFp == nil {
		path := mon.Path()
		mon.ReadFp, err = os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			return
		}
	}

	// Caller wants record 1, which stands at location 0, etc.
	_, err = mon.ReadFp.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	err = binary.Read(mon.ReadFp, binary.LittleEndian, apps)
	if err != nil {
		return
	}
	return
}

func (mon *Monitor) AppendApps(apps []index.AppearanceRecord) (count int, err error) {
	if mon.ReadFp != nil {
		mon.ReadFp.Close()
		mon.ReadFp = nil
	}

	path := mon.Path()
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	b := make([]byte, 4)
	for _, app := range apps {
		binary.LittleEndian.PutUint32(b, app.BlockNumber)
		_, err = f.Write(b)
		if err != nil {
			return
		}
		binary.LittleEndian.PutUint32(b, app.TransactionId)
		_, err = f.Write(b)
		if err != nil {
			return
		}
	}

	mon.Reload()
	count = int(mon.Count)

	return
}

func (mon *Monitor) Delete() (prev bool) {
	prev = mon.Deleted
	mon.Deleted = true
	return
}

func (mon *Monitor) UnDelete() (prev bool) {
	prev = mon.Deleted
	mon.Deleted = false
	return
}

func (mon *Monitor) ToggleDelete() (prev bool) {
	prev = mon.Deleted
	mon.Deleted = !mon.Deleted
	return
}

func (mon *Monitor) Remove() (bool, error) {
	if !mon.Deleted {
		return false, errors.New("cannot remove a file that has not been deleted")
	}
	file.Remove(mon.Path())
	return !file.FileExists(mon.Path()), nil
}
