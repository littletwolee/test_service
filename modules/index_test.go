package modules

import (
	"fmt"
	"gate/models"
	"gate/util"
	"sort"
	"testing"
	"time"

	"github.com/idoall/TokenExchangeCommon/commonmodels"
	"github.com/idoall/TokenExchangeCommon/commonstock"
)

func init() {
	util.ConfigInit("../conf")
	//util.LoggerInit()
	util.MysqlInit()
}
func Test_getIndexFromGate(t *testing.T) {
	i := NewIndex()
	_, err := i.getIndexFromGate("btc_usdt")
	if err != nil {
		t.Fatal(err)
	}
}
func Test_Sync(t *testing.T) {
	i := NewIndex()
	m := make(map[string]*models.Currency)
	m["btc_ustd"] = &models.Currency{Name: "btc_usdt"}
	i.Sync(m)
}
func Test_KDJ(t *testing.T) {
	i := NewIndex()
	list, err := i.getIndexFromGate("btc_usdt")
	if err != nil {
		t.Fatal(err)
	}
	sort.Sort(list)
	l := []*commonmodels.Kline{}
	for _, v := range list {
		t, _ := time.Parse("2006-01-02 15:04:05", time.Unix(int64(v.Dt), 0).Format("2006-01-02 15:04:05"))
		l = append(l, &commonmodels.Kline{
			Open:      v.Start,
			Close:     v.End,
			High:      v.High,
			Low:       v.Low,
			KlineTime: t,
		})
	}
	fmt.Println(len(l))
	s := commonstock.NewKDJ(l, 9)
	jdk := s.Calculation()
	for _, v := range jdk.GetPoints() {
		fmt.Printf("t:%s,k:%.8f,d:%.8f,j:%.8f\n", v.Time, v.K, v.D, v.J)
	}
}
