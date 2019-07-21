package util

import (
	"github.com/littletwolee/commons"
)

var config *conf

func Config() *conf {
	return config
}

func ConfigInit(path string) {
	v := commons.GetConfig()
	c := newConf()
	if err := v.Unmarshal(c); err != nil {
		panic(err)
	}
	config = c
}
func newConf() *conf {
	return &conf{}
}

type conf struct {
	App      app
	Services services
	Mysql    mysql
}

type app struct {
	Log log
}

type log struct {
	Path string
}

type services struct {
	Crawl crawl
}

type crawl struct {
	Currency currency
}
type currency struct {
	Interval int
	Chunk    int
}

type mysql struct {
	Host   string
	Port   int
	DBName string
	User   string
	Pwd    string
}
