package models

import "test_services/util"

type Index struct {
	ID     int64   `xorm:"id"`
	Start  float64 `xorm:"start"`
	End    float64 `xorm:"end"`
	High   float64 `xorm:"high"`
	Low    float64 `xorm:"low"`
	MA5    float64 `xorm:"ma5"`
	MA10   float64 `xorm:"ma10"`
	MA30   float64 `xorm:"ma30"`
	Dt     int     `xorm:"dt"`
	CName  string  `xorm:"cname"`
	Modefy int64   `xorm:"last_modify"`
}

const _TABLE_INDEX = "index"

func (i *Index) Insert() error {
	_, err := util.GetSession().Table(_TABLE_INDEX).Insert(i)
	return err
}
func (i *Index) Update() error {
	_, err := util.GetSession().Table(_TABLE_INDEX).Where("cname = ? and dt = ?", i.CName, i.Dt).Update(i)
	return err
}
func (i *Index) Count() (int64, error) {
	return util.GetSession().Table(_TABLE_INDEX).Where("cname = ? and dt = ?", i.CName, i.Dt).FindAndCount(i)
}
func (i *Index) Delete() error {
	_, err := util.GetSession().Sql("delete from ? where cname = ? and dt = ?", _TABLE_INDEX, i.CName, i.Dt).Exec()
	return err
}
func CountByWhere(cname string) (int64, error) {
	var i Index
	return util.GetSession().Table(_TABLE_INDEX).Where("cname = ?", cname).FindAndCount(&i)
}
func DeleteIndexByWhere(cname string, dt int) error {
	_, err := util.GetSession().Sql("delete from ? where cname = ? and dt < ?", _TABLE_INDEX, cname, dt).Exec()
	return err
}
