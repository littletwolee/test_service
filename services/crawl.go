package services

import (
	"fmt"
	"gate/modules"
	"gate/util"
	"sync"
	"time"
)

type crawl struct {
	switchChan chan bool
	wg         *sync.WaitGroup
}

func NewCrawl(wg *sync.WaitGroup) Iservice {
	return &crawl{
		switchChan: make(chan bool),
		wg:         wg,
	}
}

func (c *crawl) Start() {
	currency := modules.NewCurrency()
	index := modules.NewIndex()
	for {
		start := time.Now().UTC()
		fmt.Printf("sync index start at %s\n", start)
		next := util.Next(start, time.Duration(util.Config().Services.Crawl.Currency.Interval)*time.Second).Add(10 * time.Second)

		currencies := currency.Sync()
		index.Sync(currencies)
		end := time.Now().UTC()
		fmt.Printf("sync index completed at %s\n", end)
		if end.After(next) {
			continue
		}
		time.Sleep(util.Interval(next, end))
	}
	// p := pipe.NewPipe(c.conf.Services.Crawl.Chunk)
	// p.AddJobs()
}
func (c *crawl) Stop() {
}
