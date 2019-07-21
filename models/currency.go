package models

import (
	"fmt"
	"gate/util"
	"strconv"
	"strings"
	"time"
)

type Currency struct {
	Name    string  `xorm:"name" `
	Volume  float64 `xorm:"volume"`
	Modefy  int64   `xorm:"last_modify"`
	DBEvent DBEvent `xorm:"-"`
}

const _TABLE_CURRENCY = "currency"

func (c *Currency) Insert() error {
	_, err := util.GetSession().Table(_TABLE_CURRENCY).Insert(c)
	return err
}
func (c *Currency) Update() error {
	_, err := util.GetSession().Exec("update "+_TABLE_CURRENCY+" set volume=?,last_modify=? where name=?;", c.Volume, c.Modefy, c.Name)
	return err
}
func (c *Currency) Exists() (bool, error) {
	return util.GetSession().Table(_TABLE_CURRENCY).Exist(c)
}
func (c *Currency) Delete() error {
	_, err := util.GetSession().Exec("delete from "+_TABLE_CURRENCY+" where name=?;", c.Name)
	return err
}
func AllIDs(list *[]*Currency) error {
	return util.GetSession().Table(_TABLE_CURRENCY).Select("name").Find(list)
}
func All(list *[]*Currency) error {
	return util.GetSession().Table(_TABLE_CURRENCY).Find(list)
}

type JsonCurrency struct {
	ID     int    `json:"no"`
	Name   string `json:"curr_a"`
	Volume string `json:"vol_b"`
	Suffix string `json:"curr_b"`
}

func (jc *JsonCurrency) ToCurrency(event DBEvent) *Currency {
	jc.Volume = strings.Replace(jc.Volume, ",", "", -1)
	vol, err := strconv.ParseFloat(jc.Volume, 64)
	if err != nil {
		fmt.Printf("parse field volume error: %s\n", err.Error())
	}
	return &Currency{
		Name:    jc.Name,
		Volume:  vol,
		DBEvent: event,
		Modefy:  time.Now().UTC().Unix(),
	}
}

func (jc *JsonCurrency) FormatName() {
	jc.Name = strings.ToLower(jc.Name) + "_" + strings.ToLower(jc.Suffix)
}
