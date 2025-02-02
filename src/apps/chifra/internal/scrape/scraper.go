package scrapePkg

// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

import (
	"io/ioutil"
	"log"
	"time"
)

var statusPath string = "/tmp/"

type Scraper struct {
	Running   bool    `json:"Running"`
	SleepSecs float64 `json:"SleepSecs"`
	Name      string  `json:"Name"`
}

func NewScraper(color, name string, secs float64, logLev uint64) Scraper {
	scraper := new(Scraper)
	scraper.Name = name
	scraper.SleepSecs = secs
	scraper.Running = false
	return *scraper
}

func (scraper *Scraper) ChangeState(onOff bool) bool {
	prev := scraper.Running
	scraper.Running = onOff
	str := "false"
	if onOff {
		str = "true"
	}
	fileName := statusPath + scraper.Name + ".txt"
	err := ioutil.WriteFile(fileName, []byte(str), 0644)
	if err != nil {
		log.Fatal(err)
	}
	return prev
}

func (scraper *Scraper) Pause() {
	halfSecs := scraper.SleepSecs * 2
	state := scraper.Running
	for i := 0; i < int(halfSecs); i++ {
		if state != scraper.Running {
			break
		}
		time.Sleep(time.Duration(500) * time.Millisecond)
	}
}
