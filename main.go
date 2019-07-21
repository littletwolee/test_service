package main

import (
	"fmt"
	"gate/services"
	"gate/util"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	_init()
	wg := new(sync.WaitGroup)
	servicesList := initServices(wg)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	start(servicesList)
	stop(wg, ch, servicesList)
}
func _init() {
	util.ConfigInit("./conf")
	//util.LoggerInit()
	util.MysqlInit()
}

func initServices(wg *sync.WaitGroup) []services.Iservice {
	return []services.Iservice{
		services.NewCrawl(wg),
		services.NewServer(wg),
	}
}

func start(servicesList []services.Iservice) {
	fmt.Println("app start")
	for _, v := range servicesList {
		go v.Start()
	}
}
func stop(wg *sync.WaitGroup, exit chan os.Signal, servicesList []services.Iservice) {
	<-exit
	// wg.Add(len(servicesList))
	for _, v := range servicesList {
		v.Stop()
	}
	// wg.Wait()
	fmt.Println("app shutdown")
}
