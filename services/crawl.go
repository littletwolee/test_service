package services

import (
	"sync"
	"test_services/modules"
	"test_services/util"
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
	for {
		start := time.Now().UTC()
		util.Logger().InfoF("sync currency start at %s", start)
		next := util.Next(start, time.Duration(util.Config().Services.Crawl.Currency.Interval)*time.Second).Add(10 * time.Second)

		currency.Sync()
		end := time.Now().UTC()
		util.Logger().InfoF("sync currency completed at %s", end)
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
