package modules

import (
	"test_services/util"
	"testing"
)

func init() {
	util.ConfigInit("../conf")
	util.LoggerInit()
}
func Test_getCurrencies(t *testing.T) {
	// c := NewCurrency()
	// if err := c.getCurrenciesFromGate(); err != nil {
	// 	t.Fatal(err)
	// }
}
