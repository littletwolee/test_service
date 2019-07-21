package models

import (
	"encoding/json"
	"gate/util"
	"time"
)

type IndexDB struct {
	T             time.Time `json:"t",xorm:"t"`
	CurrenciesStr string    `json:"-" xorm:"currencies"`
	Currencies    []string  `json:"currencies" xorm:"-"`
}

const _TABLE_INDEXS = "indexs"

func (i *IndexDB) Insert() error {
	_, err := util.GetSession().Table(_TABLE_INDEXS).Insert(i)
	return err
}

func (i *IndexDB) Delete() error {
	var ii IndexDB
	_, err := util.GetSession().Table(_TABLE_INDEXS).Where("t < ?", i.T).Delete(&ii)
	return err
}

type IndexDBs []*IndexDB

func (i *IndexDBs) ToList() {
	for _, v := range *i {
		var c []string
		json.Unmarshal([]byte(v.CurrenciesStr), &c)
		v.Currencies = c
	}
}

func (i *IndexDBs) Select() error {
	return util.GetSession().Table(_TABLE_INDEXS).Desc("t").Find(i)
}
