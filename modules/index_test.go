package modules

import (
	"test_services/models"
	"test_services/util"
	"testing"
)

func init() {
	util.ConfigInit("../conf")
	util.LoggerInit()
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
