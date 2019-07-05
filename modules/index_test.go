package modules

import (
	"test_services/util"
	"testing"
)

func init() {
	util.ConfigInit("../conf")
	util.LoggerInit()
}
func Test_getIndexFromGate(t *testing.T) {
	i := NewIndex()
	_, err := i.getIndexFromGate("btc_usdt")
	if err != nil {
		t.Fatal(err)
	}
}
