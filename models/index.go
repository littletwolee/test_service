package models

import "test_services/util"

type Index struct {
	Start  float64 `xorm:"start"`
	End    float64 `xorm:"end"`
	High   float64 `xorm:"high"`
	Low    float64 `xorm:"low"`
	MA5    float64 `xorm:"ma5"`
	MA10   float64 `xorm:"ma10"`
	MA30   float64 `xorm:"ma30"`
	K      float64 `xorm:"k"`
	D      float64 `xorm:"d"`
	J      float64 `xorm:"j"`
	Dt     int     `xorm:"dt"`
	CName  string  `xorm:"cname"`
	Modefy int64   `xorm:"last_modify"`
}

type Indexs []*Index

func (is Indexs) Len() int { return len(is) }

func (is Indexs) Swap(i, j int) { is[i], is[j] = is[j], is[i] }

func (is Indexs) Less(i, j int) bool { return is[i].Dt < is[j].Dt }

func (is Indexs) MA(i, l int) float64 {
	var result float64
	max := is.Len()
	if max >= i+l {
		max = i + l
	}
	for k := i; k < max; k++ {
		result += is[k].End
	}
	return result / float64(max-i)
}

// func (is Indexs) K() float64 {

// }
func (is Indexs) rsv(i, l int) float64 {
	low, high := is.lowAndHigh(i, l)
	return (is[i].End - low) / (high - low) * 100
}

func (is Indexs) lowAndHigh(i, l int) (float64, float64) {
	var high, low float64
	max := is.Len()
	if max >= i+l {
		max = i + l
	}
	low = is[i].Low
	high = is[i].High
	for k := i; k < max; k++ {
		if is[k].Low < low {
			low = is[k].Low
		}
		if is[k].High > high {
			high = is[k].High
		}
	}
	return low, high
}

const _TABLE_INDEX = "indexs"

func (i *Index) Insert() error {
	_, err := util.GetSession().Table(_TABLE_INDEX).Insert(i)
	return err
}
func (i *Index) Update() error {
	_, err := util.GetSession().Table(_TABLE_INDEX).Where("cname = ? and dt = ?", i.CName, i.Dt).Update(i)
	return err
}
func (i *Index) Count() (int64, error) {
	return util.GetSession().Table(_TABLE_INDEX).Where("cname = ? and dt = ?", i.CName, i.Dt).Count()
}
func (i *Index) Delete() error {
	_, err := util.GetSession().Sql("delete from ? where cname = ? and dt = ?", _TABLE_INDEX, i.CName, i.Dt).Exec()
	return err
}
func CountByWhere(cname string) (int64, error) {
	var is []Index
	return util.GetSession().Table(_TABLE_INDEX).Where("cname = ?", cname).FindAndCount(&is)
}
func DeleteIndexByWhere(cname string, dt int) error {
	var i Index
	_, err := util.GetSession().Table(_TABLE_INDEX).Where("cname = ? and dt < ?", cname, dt).Delete(&i)
	return err
}
