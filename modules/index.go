package modules

import (
	"encoding/json"
	"fmt"
	"gate/models"
	"gate/util"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type index struct{ host string }

type names struct {
	m    *sync.RWMutex
	t    time.Time
	list []string
}

func (n *names) add(l []string) {
	n.m.Lock()
	defer n.m.Unlock()
	n.list = append(n.list, l...)
}

func NewIndex() *index {
	return &index{
		host: "https://www.gateio.io/json_svr/query",
	}
}
func (i *index) Sync(list []*models.JsonCurrency) {
	ns := &names{m: new(sync.RWMutex), list: []string{}}
	m := len(list) / 10
	n := len(list) % 10
	s := 0
	if n != 0 {
		m++
	}
	wg := new(sync.WaitGroup)
	wg.Add(10)
	for j := 0; j < 10; j++ {
		if j == 0 {
			go i.getData(list[s:m], ns, wg, true)
		} else {
			go i.getData(list[s:m], ns, wg, false)
		}
		s = m
		m += m
		if m > len(list) {
			m = len(list)
		}
	}
	wg.Wait()
	i.saveToDB(ns)
}

func (i *index) getData(vs []*models.JsonCurrency, ns *names, wg *sync.WaitGroup, flag bool) {
	var list []string
	for _, v := range vs {
		data, err := i.getIndexFromGate(v.Name)
		if err != nil {
			fmt.Printf("get %s data from gate error: %s\n", v.Name, err.Error())
		}
		calculateMA(&data)
		if data != nil && len(data) >= 5 {
			d := data[1]
			if flag {
				ns.t = d.Dt.T()
			}
			flag = false
			d1 := data[2]
			d2 := data[3]
			d3 := data[4]
			if d.MA5 >= d.MA10 && d.MA10 >= d.MA30 {
				if (d1.MA30 >= d1.MA10 && d1.MA30 >= d1.MA5) || (d2.MA30 >= d2.MA10 && d2.MA30 >= d3.MA5) || (d3.MA30 >= d3.MA10 && d3.MA30 >= d3.MA5) {
					list = append(list, d.CName)
				}
			}
		}
	}
	ns.add(list)
	wg.Done()
}

func (i *index) getIndexFromGate(cname string) (models.Indexs, error) {
	cli := util.NewClient()
	n := time.Now()
	to := n.Unix()
	from := n.Add(-20 * time.Hour).Unix()
	err := cli.ParseUrl(i.host).Query("u", "10").Query("c", "4846799").Query("type", "tvkline").Query("symbol", cname).Query("interval", "1800").Query("from", strconv.Itoa(int(from))).Query("to", strconv.Itoa(int(to))).Do(http.MethodGet)
	if err != nil {
		return nil, err
	}
	var list models.Indexs
	//fmt.Println(cli.Url())
	for _, v := range strings.Split(string(cli.Body()), "\n")[1:] {
		data := strings.Split(v, ",")
		if len(data) < 5 {
			continue
		}
		list = append(list, &models.Index{
			CName: cname,
			Dt:    models.TimeInt(i.toFInt(data[0])),
			Start: i.toFloat(data[1]),
			High:  i.toFloat(data[2]),
			Low:   i.toFloat(data[3]),
			End:   i.toFloat(data[4]),
		})
	}
	sort.Sort(sort.Reverse(list))
	return list, nil
}

func calculateMA(is *models.Indexs) {
	for i, v := range *is {
		v.MA5 = is.MA(i, 5)
		v.MA10 = is.MA(i, 10)
		v.MA30 = is.MA(i, 30)
	}
}

func (i *index) saveToDB(names *names) {
	str, err := json.Marshal(names.list)
	if err != nil {
		fmt.Printf("marshal index list error: %s\n", err.Error())
	}
	index := &models.IndexDB{T: names.t, CurrenciesStr: string(str)}
	if err := index.Insert(); err != nil {
		fmt.Printf("insert indexdb error:%s\n", err.Error())
	}
	index = &models.IndexDB{T: names.t.Add(-12 * time.Hour)}
	if err := index.Delete(); err != nil {
		fmt.Printf("delete indexdb error:%s\n", err.Error())
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
