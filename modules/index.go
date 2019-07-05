package modules

import (
	"net/http"
	"strconv"
	"strings"
	"test_services/models"
	"test_services/util"
)

type index struct{ host string }

func NewIndex() *index {
	return &index{
		host: "https://www.gateio.io/json_svr/query",
	}
}
func (i *index) Sync(list []*models.Currency) {
	for _, v := range list {
		data, err := i.getIndexFromGate(v.Name)
		if err != nil {
			util.Logger().ErrorF("get %s data from gate error: %s", v.Name, err.Error())
		}
		i.saveToDB(v.Name, data)
	}
}

func (i *index) getIndexFromGate(cname string) ([]*models.Index, error) {
	cli := util.NewClient()
	err := cli.ParseUrl(i.host).Query("u", "10").Query("c", "4846799").Query("type", "kline").Query("symbol", cname).Query("group_sec", "1800").Query("range_hour", "4.7").Do(http.MethodGet)
	if err != nil {
		return nil, err
	}
	var list []*models.Index
	for _, v := range strings.Split(string(cli.Body()), "\n")[1:] {
		data := strings.Split(v, ",")
		list = append(list, &models.Index{
			Dt:    i.toFInt(data[0]),
			Start: i.toFloat(data[1]),
			High:  i.toFloat(data[2]),
			Low:   i.toFloat(data[3]),
			End:   i.toFloat(data[4]),
		})
	}
	return list, nil
}

func (i *index) saveToDB(cname string, list []*models.Index) {
	for _, v := range list {
		count, err := v.Count()
		if err != nil {
			util.Logger().ErrorF("search cname: %s, dt: %d, exists error: %s", cname, v.Dt, err.Error())
		}
		if count > 0 {
			if err := v.Update(); err != nil {
				util.Logger().ErrorF("update index error: %s, data: %+v", err.Error(), v)
			}
		} else {
			if err := v.Insert(); err != nil {
				util.Logger().ErrorF("insert index error: %s, data: %+v", err.Error(), v)
			}
		}

	}
	count, err := models.CountByWhere(cname)
	if err != nil {
		util.Logger().ErrorF("search count error: %s, cname: %+v", err.Error(), cname)
		return
	}
	if count > 48 {
		if err := models.DeleteIndexByWhere(cname, list[0].Dt); err != nil {
			util.Logger().ErrorF("delete dirty index error: %s, cname: %+v", err.Error(), cname)
		}
	}

}

func (i *index) toFloat(v string) float64 {
	f, _ := strconv.ParseFloat(v, 64)
	return f
}
func (i *index) toFInt(v string) int {
	d, _ := strconv.Atoi(v)
	return d / 1000
}
